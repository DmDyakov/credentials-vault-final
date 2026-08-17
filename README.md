# Credentials Vault

Менеджер приватных данных: логины, пароли, банковские карты,
текстовые заметки и бинарные файлы.

## Быстрый старт

```bash
make prod          # запуск в Docker
make dev           # локальная разработка
make build-all     # сборка всех бинарников
```

## Тестирование

```bash
make check         # все проверки
make test          # тесты
make test-coverage # покрытие
```

## API

gRPC сервисы описаны в `api/proto/`:

- `auth/v1/auth.proto` — регистрация и вход
- `vault/v1/vault.proto` — хранение данных

## Структура

Монорепо с Go workspace:

- `server/` — gRPC сервер
- `client-cli/` — CLI клиент
- `gen/` — сгенерированный код
- `pkg/` — общие пакеты
- `tools/` — инструменты разработки
