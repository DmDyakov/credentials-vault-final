VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -X credentials-vault/pkg/buildinfo.Version=$(VERSION) \
           -X credentials-vault/pkg/buildinfo.Date=$(DATE) \
           -X credentials-vault/pkg/buildinfo.Commit=$(COMMIT) \
           -s -w

MODULES := server client-cli gen pkg tools/staticlint

BIN_EXT := $(shell go env GOEXE)

.PHONY: proto mocks generate setup build build-cli build-all dev prod down test test-coverage test-coverage-html test-race \
        check-fmt fmt lint vet staticlint check \
        migrate-up migrate-down migrate-status docker-logs docker-logs-db docker-ps docker-rebuild \
        clean clean-all

# ═══════════════════════════════════════════
# Генерация и сборка
# ═══════════════════════════════════════════

# Генерация gRPC-кода (Edition 2024)
# Генерация gRPC-кода
proto:
	mkdir -p gen/go/auth/v1 gen/go/vault/v1
	protoc \
		-I api/proto \
		--go_out=. \
		--go_opt=module=credentials-vault \
		--go-grpc_out=. \
		--go-grpc_opt=module=credentials-vault \
		api/proto/auth/v1/auth.proto \
		api/proto/vault/v1/vault.proto


# Генерация моков
mocks:
	go generate -C server ./...
	go generate -C client-cli ./...

# Генерация всего (proto + mocks)
generate: proto mocks

# Полная настройка: зависимости, генерация, сборка
setup: generate
	go work sync
	make build-all

# Сборка сервера
build:
	mkdir -p bin
	go build -C server -trimpath -ldflags "$(LDFLAGS)" -o ../bin/server$(BIN_EXT) ./cmd/server

# Сборка CLI
build-cli:
	mkdir -p bin
	go build -C client-cli -trimpath -ldflags "$(LDFLAGS)" -o ../bin/vault$(BIN_EXT) ./cmd/cli

# Сборка всех бинарников
build-all: build build-cli
	@echo "✅ All binaries built!"

# ═══════════════════════════════════════════
# Окружения
# ═══════════════════════════════════════════

# Dev-окружение
dev:
	docker compose --env-file .env.dev up -d --no-deps postgres
	@bash -c 'set -a && source .env.dev && set +a && go run -C server ./cmd/server'

# Prod-окружение: всё в Docker
prod:
	docker compose --env-file .env.prod up -d

# Остановить все контейнеры
down:
	docker compose down

# Остановить и удалить volumes
down-clean:
	docker compose down -v

# ═══════════════════════════════════════════
# Тестирование
# ═══════════════════════════════════════════

# Тесты
test:
	@for mod in $(MODULES); do \
		echo "=== test $$mod ==="; \
		go test -C $$mod ./... || exit 1; \
	done

# Тесты с покрытием
test-coverage:
	@for mod in $(MODULES); do \
		echo "=== test-coverage $$mod ==="; \
		coverage_file="coverage_$$(echo $$mod | tr '/' '_').out"; \
		go test -C $$mod -coverprofile=$$coverage_file -covermode=atomic ./... || exit 1; \
	done
	@echo "Coverage files: coverage_*.out"

# HTML отчет покрытия
test-coverage-html:
	@make test-coverage
	@go tool cover -html=coverage_server.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Тесты с race detector
test-race:
	@for mod in $(MODULES); do \
		echo "=== test-race $$mod ==="; \
		go test -race -C $$mod ./... || exit 1; \
	done

# ═══════════════════════════════════════════
# Линтеры
# ═══════════════════════════════════════════

# Форматирование
fmt:
	@for mod in $(MODULES); do \
		echo "=== fmt $$mod ==="; \
		gofmt -w $$mod; \
	done

# Проверка форматирования
check-fmt:
	@for mod in $(MODULES); do \
		unformatted=$$(gofmt -l $$mod); \
		if [ -n "$$unformatted" ]; then \
			echo "Files need formatting in $$mod:"; \
			echo "$$unformatted"; \
			exit 1; \
		fi; \
	done

# Стандартный линтер
lint:
	@for mod in $(MODULES); do \
		echo "=== lint $$mod ==="; \
		cd $$mod && golangci-lint run ./... && cd .. || exit 1; \
	done

# Go vet
vet:
	@for mod in $(MODULES); do \
		echo "=== vet $$mod ==="; \
		go vet -C $$mod ./... || exit 1; \
	done

# Кастомный статический анализатор
staticlint:
	go run ./tools/staticlint/cmd/staticlint ./server/... ./client-cli/... ./gen/... ./pkg/...

# Все проверки
check: check-fmt vet lint staticlint test
	@echo "✅ All checks passed!"

# ═══════════════════════════════════════════
# Миграции
# ═══════════════════════════════════════════

# Применить миграции
migrate-up:
	docker compose exec postgres psql -U postgres -d credentials_vault -c "SELECT 1" > /dev/null 2>&1 || docker compose up -d postgres
	sleep 2
	docker compose run --rm server ./server -migrate

# Откатить миграции
migrate-down:
	docker compose run --rm server ./server -migrate-down

# Статус миграций
migrate-status:
	docker compose exec postgres psql -U postgres -d credentials_vault -c "SELECT * FROM schema_migrations;"

# ═══════════════════════════════════════════
# Docker
# ═══════════════════════════════════════════

# Логи сервера
docker-logs:
	docker compose logs -f server

# Логи БД
docker-logs-db:
	docker compose logs -f postgres

# Статус контейнеров
docker-ps:
	docker compose ps

# Пересборка
docker-rebuild:
	docker compose up -d --build

# ═══════════════════════════════════════════
# Очистка
# ═══════════════════════════════════════════

# Очистка артефактов
clean:
	rm -rf bin/
	rm -f coverage_*.out coverage.html
	rm -rf gen/auth gen/vault

# Полная очистка
clean-all: clean down-clean
	@echo "✅ Cleaned everything!"
