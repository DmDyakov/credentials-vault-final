VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -X credentials-vault/pkg/buildinfo.Version=$(VERSION) \
           -X credentials-vault/pkg/buildinfo.Date=$(DATE) \
           -X credentials-vault/pkg/buildinfo.Commit=$(COMMIT) \
           -s -w

MODULES := server client-cli gen pkg tools/staticlint

BIN_EXT := $(shell go env GOEXE)

.PHONY: proto mocks generate setup build build-cli build-all dev prod down test test-integration test-all \
        test-coverage test-coverage-html test-race \
        check-fmt fmt lint vet staticlint check \
        migrate-up migrate-down migrate-status docker-logs docker-logs-db docker-ps docker-rebuild \
        clean clean-all

# ═══════════════════════════════════════════
# Генерация и сборка
# ═══════════════════════════════════════════

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
		@find gen -name "*.pb.go" -exec bash -c 'mv "$$0" "$${0%.pb.go}.pb.gen.go"' {} \;

mocks:
	go generate -C server ./...
	go generate -C client-cli ./...

generate: proto mocks

setup: generate
	go work sync
	make build-all

build:
	mkdir -p bin
	go build -C server -trimpath -ldflags "$(LDFLAGS)" -o ../bin/server$(BIN_EXT) ./cmd/server

build-cli:
	mkdir -p bin
	go build -C client-cli -trimpath -ldflags "$(LDFLAGS)" -o ../bin/vault$(BIN_EXT) ./cmd/cli

build-all: build build-cli
	@echo "✅ All binaries built!"

# ═══════════════════════════════════════════
# Окружения
# ═══════════════════════════════════════════

dev:
	docker compose --env-file .env.dev up -d --no-deps postgres
	@bash -c 'set -a && source .env.dev && set +a && go run ./server/cmd/server'

prod:
	docker compose --env-file .env.prod up -d

down:
	docker compose down

down-clean:
	docker compose down -v

# ═══════════════════════════════════════════
# Тестирование
# ═══════════════════════════════════════════

test:
	@for mod in $(MODULES); do \
		echo "=== test $$mod ==="; \
		go test -C $$mod ./... || exit 1; \
	done

test-integration:
	cd client-cli && go test -tags=integration -run TestClientFullFlow -v ./test/integration/

test-all: test test-integration
	@echo "✅ All tests passed!"

test-coverage:
	@for mod in $(MODULES); do \
		echo "=== test-coverage $$mod ==="; \
		coverage_file="coverage_$$(echo $$mod | tr '/' '_').out"; \
		go test -C $$mod -coverprofile=$$coverage_file -covermode=atomic ./... || exit 1; \
	done
	@echo "Coverage files: coverage_*.out"

test-coverage-html:
	@make test-coverage
	@go tool cover -html=coverage_server.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-race:
	@for mod in $(MODULES); do \
		echo "=== test-race $$mod ==="; \
		go test -race -C $$mod ./... || exit 1; \
	done

# ═══════════════════════════════════════════
# Линтеры
# ═══════════════════════════════════════════

fmt:
	@for mod in $(MODULES); do \
		echo "=== fmt $$mod ==="; \
		gofmt -w $$mod; \
	done

check-fmt:
	@for mod in $(MODULES); do \
		unformatted=$$(gofmt -l $$mod); \
		if [ -n "$$unformatted" ]; then \
			echo "Files need formatting in $$mod:"; \
			echo "$$unformatted"; \
			exit 1; \
		fi; \
	done

lint:
	@for mod in $(MODULES); do \
		echo "=== lint $$mod ==="; \
		cd $$mod && golangci-lint run ./... && cd .. || exit 1; \
	done

vet:
	@for mod in $(MODULES); do \
		echo "=== vet $$mod ==="; \
		go vet -C $$mod ./... || exit 1; \
	done

staticlint:
	go run ./tools/staticlint/cmd/staticlint ./server/... ./client-cli/... ./gen/... ./pkg/...

check: check-fmt vet lint staticlint test test-integration
	@echo "✅ All checks passed!"

# ═══════════════════════════════════════════
# Миграции
# ═══════════════════════════════════════════

migrate-up:
	docker compose exec postgres psql -U postgres -d credentials_vault -c "SELECT 1" > /dev/null 2>&1 || docker compose up -d postgres
	sleep 2
	docker compose run --rm server ./server -migrate

migrate-down:
	docker compose run --rm server ./server -migrate-down

migrate-status:
	docker compose exec postgres psql -U postgres -d credentials_vault -c "SELECT * FROM schema_migrations;"

# ═══════════════════════════════════════════
# Docker
# ═══════════════════════════════════════════

docker-logs:
	docker compose logs -f server

docker-logs-db:
	docker compose logs -f postgres

docker-ps:
	docker compose ps

docker-rebuild:
	docker compose up -d --build

# ═══════════════════════════════════════════
# Очистка
# ═══════════════════════════════════════════

clean:
	rm -rf bin/
	rm -f coverage_*.out coverage.html
	rm -rf gen/auth gen/vault

clean-all: clean down-clean
	@echo "✅ Cleaned everything!"
