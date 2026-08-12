# Credentials Vault

Менеджер приватных данных: логины, пароли, банковские карты,
текстовые заметки и бинарные файлы.

## Быстрый старт

```bash
make prod
grpcurl -plaintext localhost:9090 list
```

## Разработка

```bash
make dev
```

## Полезные команды

```bash
make proto           # генерация gRPC-кода
make build           # сборка бинарника
make test            # тесты
make test-coverage   # покрытие
make lint            # golangci-lint
make staticlint      # кастомный анализатор
```
