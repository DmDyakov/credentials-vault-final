VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -X credentials-vault/pkg/buildinfo.Version=$(VERSION) \
           -X credentials-vault/pkg/buildinfo.Date=$(DATE) \
           -X credentials-vault/pkg/buildinfo.Commit=$(COMMIT)

.PHONY: proto mocks generate setup build dev prod down test test-coverage test-coverage-html check-fmt fmt lint staticlint \
        migrate-up migrate-down migrate-status docker-logs docker-ps docker-rebuild clean

# ═══════════════════════════════════════════
# Генерация и сборка
# ═══════════════════════════════════════════

# Генерация gRPC-кода
proto:
	mkdir -p gen/go
	protoc \
		-I api/proto \
		--go_out=gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=gen/go \
		--go-grpc_opt=paths=source_relative \
		api/proto/auth/v1/auth.proto \
		api/proto/vault/v1/vault.proto

# Генерация моков
mocks:
	go generate ./...

# Генерация всего (proto + mocks)
generate: proto mocks

# Полная настройка: зависимости, генерация, сборка
setup: generate
	go mod tidy
	make build

# Сборка бинарника
build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/server ./cmd/server

# ═══════════════════════════════════════════
# Окружения
# ═══════════════════════════════════════════

# Dev-окружение
dev:
	docker compose --env-file .env.dev up -d --no-deps postgres
	@bash -c 'set -a && source .env.dev && set +a && go run ./cmd/server'

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
	go test ./...

# Тесты с покрытием
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# HTML отчет покрытия
test-coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Тесты с race detector
test-race:
	go test -race ./...

# ═══════════════════════════════════════════
# Линтеры
# ═══════════════════════════════════════════

# Форматирование
fmt:
	go fmt ./...

# Проверка форматирования
check-fmt:
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:" && gofmt -l . && exit 1)

# Стандартный линтер
lint:
	golangci-lint run ./...

# Go vet
vet:
	go vet ./...

# Кастомный статический анализатор
staticlint:
	go run ./cmd/staticlint/ ./...

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
	rm -f coverage.out coverage.html
	rm -rf gen/go/

# Полная очистка
clean-all: clean down-clean
	@echo "✅ Cleaned everything!"