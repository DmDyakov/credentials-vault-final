# Credentials Vault

Credentials Vault — система безопасного хранения приватных данных: логинов, паролей, банковских карт, текстовых заметок и бинарных файлов.

## Возможности

- **End-to-end шифрование** — данные шифруются на клиенте, сервер не имеет доступа к содержимому
- **Кроссплатформенный CLI** — Windows, Linux, macOS
- **Синхронизация** — доступ к данным с любого устройства
- **Метаинформация** — произвольные текстовые метки для любых данных

## Безопасность

### Архитектура

Сервер не получает мастер-пароль или ключи шифрования. Данные шифруются на клиенте до отправки на сервер.

### Аутентификация

- Пользователь вводит мастер-пароль
- Сервер хранит bcrypt-хеш пароля
- После входа выдаётся JWT-токен
- Токен используется для авторизации запросов

### Шифрование данных

- Ключ шифрования: `Argon2id(master_password, salt)` — 256 бит
- Шифрование: `AES-256-GCM`
- Соль генерируется при регистрации и хранится на сервере
- Шифрование и расшифровка выполняются исключительно на клиенте

### Защита канала

- Канал связи защищён TLS
- Двухуровневое шифрование: данные шифруются до отправки, затем защищаются TLS

### Восстановление доступа

При смене устройства пользователю достаточно ввести мастер-пароль. Соль загружается с сервера, ключ шифрования восстанавливается, данные расшифровываются.

## Установка

### Сборка из исходников

```bash
make build-all
```

Бинарники будут размещены в `bin/`.

### Запуск сервера

```bash
docker compose up -d
```

## Использование

### Системные команды

```bash
# Инициализация конфигурации
vault init --server=localhost:9090 --ca-file=certs/server.crt

# Версия
vault version
```

### Команды аутентификации

```bash
# Регистрация (пароль вводится скрыто)
vault register --username=user

# Вход (пароль вводится скрыто)
vault login --username=user

# Выход
vault logout
```

### Команды работы с данными

```bash
# Добавление учётных данных (пароль вводится скрыто)
vault add credentials --site=example.com --username=user

# Добавление банковской карты (CVV вводится скрыто)
vault add card --brand=visa --bank=sberbank --number=4111111111111111 --holder="IVAN IVANOV" --expiry=12/25

# Список записей
vault list

# Получение записи
vault get <id>
```

### Windows (Git Bash / WSL)

```bash
# Инициализация
./bin/vault.exe init --server=localhost:9090 --ca-file=certs/server.crt

# Регистрация
./bin/vault.exe register --username=user

# Вход
./bin/vault.exe login --username=user

# Добавление учётных данных
./bin/vault.exe add credentials --site=example.com --username=user

# Добавление карты
./bin/vault.exe add card --brand=visa --bank=sberbank --number=4111111111111111 --holder="IVAN IVANOV" --expiry=12/25

# Список записей
./bin/vault.exe list

# Получение записи
./bin/vault.exe get <id>

# Версия
./bin/vault.exe version
```

### Windows (PowerShell / CMD)

```powershell
# Инициализация
.\bin\vault.exe init --server=localhost:9090 --ca-file=certs/server.crt

# Регистрация
.\bin\vault.exe register --username=user

# Вход
.\bin\vault.exe login --username=user

# Добавление учётных данных
.\bin\vault.exe add credentials --site=example.com --username=user

# Добавление карты
.\bin\vault.exe add card --brand=visa --bank=sberbank --number=4111111111111111 --holder="IVAN IVANOV" --expiry=12/25

# Список записей
.\bin\vault.exe list

# Получение записи
.\bin\vault.exe get <id>

# Версия
.\bin\vault.exe version
```

## Разработка

```bash
make dev              # Запуск в dev-режиме
make test             # Запуск unit-тестов
make test-integration # Запуск интеграционных тестов
make test-all         # Все тесты
make test-coverage    # Покрытие кода
make check            # Линтеры и проверки
make build            # Сборка сервера
make build-cli        # Сборка CLI
make build-all        # Сборка всех бинарников
```

## Технологии

- Go 1.26
- gRPC + Protobuf
- PostgreSQL (pgxpool)
- Argon2id
- AES-256-GCM
- TLS
- Docker
- Cobra

## Лицензия

MIT
