VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -X credentials-vault/pkg/buildinfo.Version=$(VERSION) \
           -X credentials-vault/pkg/buildinfo.Date=$(DATE) \
           -X credentials-vault/pkg/buildinfo.Commit=$(COMMIT)

.PHONY: proto setup build dev prod down test test-coverage check-fmt lint staticlint

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

# Полная настройка: зависимости, генерация, сборка
setup: proto
	go mod tidy
	make build

# Сборка бинарника
build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/server ./cmd/server

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

# Тесты
test:
	go test ./...

# Тесты с покрытием
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Проверка форматирования
check-fmt:
	test -z "$$(gofmt -l .)"

# Стандартный линтер
lint:
	golangci-lint run ./...

# Кастомный статический анализатор
staticlint:
	go run ./cmd/staticlint/ ./...