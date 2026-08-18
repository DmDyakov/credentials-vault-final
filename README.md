# Credentials Vault

Менеджер приватных данных: логины, пароли, банковские карты,
текстовые заметки и бинарные файлы.

## Быстрый старт

```bash
# Запуск сервера
docker compose up -d

# Сборка CLI
make build-cli

# Регистрация
./bin/vault register --username=user --password=pass

# Вход
./bin/vault login --username=user --password=pass

# Добавление логина
./bin/vault add login --site=example.com --username=user --password=pass

# Список элементов
./bin/vault list

# Получение элемента
./bin/vault get <id>
```

> **Примечание:** На Windows используйте `vault.exe` вместо `vault`.

## Команды CLI

```bash
vault register              # Регистрация нового пользователя
vault login                 # Вход в систему
vault add login             # Добавить логин/пароль
vault list                  # Список элементов
vault get <id>              # Получить элемент по ID
vault version               # Версия и дата сборки
```

## Кроссплатформенная сборка

```bash
# Linux (amd64)
make build-linux

# macOS (arm64)
make build-darwin

# Windows (amd64)
make build-windows

# Все платформы
make build-all-platforms
```

Бинарники будут в `bin/`:

## Разработка

```bash
make dev          # Запуск в dev-режиме
make test         # Тесты всех модулей
make test-coverage # Покрытие кода
make check        # Все проверки (fmt, vet, lint, staticlint, test)
make build        # Сборка сервера
make build-cli    # Сборка CLI
make build-all    # Сборка всех бинарников
```

## Структура

Монорепо с Go workspace:

- `server/` — gRPC сервер
- `client-cli/` — CLI клиент
- `gen/` — сгенерированный gRPC код
- `pkg/` — общие пакеты (buildinfo, jwt, lifecycle)
- `tools/` — инструменты разработки (staticlint)
- `api/` — proto-файлы

## API

gRPC сервисы:

- `AuthService` — регистрация и аутентификация
- `VaultService` — хранение и управление данными

## Технологии

- Go 1.26
- gRPC + Protobuf
- PostgreSQL
- JWT
- Docker
- Cobra
