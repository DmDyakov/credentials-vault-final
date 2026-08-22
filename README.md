```markdown
# Credentials Vault

Credentials Vault — система безопасного хранения приватных данных: логинов, паролей, банковских карт, текстовых заметок и бинарных файлов.

## Возможности

- **End-to-end шифрование** — данные шифруются на клиенте, сервер не имеет доступа к содержимому
- **Zero-knowledge аутентификация** — мастер-пароль никогда не покидает устройство
- **Кроссплатформенный CLI** — Windows, Linux, macOS
- **Синхронизация** — доступ к данным с любого устройства
- **Метаинформация** — произвольные текстовые метки для любых данных

## Безопасность

### Архитектура

Система построена на принципе zero-knowledge: сервер не получает мастер-пароль, ключи шифрования или расшифрованные данные.

### Аутентификация

- Пользователь вводит мастер-пароль
- На клиенте вычисляется приватный ключ: `Ed25519 = Argon2id(master_password, salt)`
- На сервер передаётся только публичный ключ
- Аутентификация выполняется по протоколу challenge-response с подписью Ed25519
- Сервер хранит публичный ключ — его нельзя использовать для входа

### Шифрование данных

- Ключ шифрования: `Argon2id(master_password, salt)` — 256 бит
- Шифрование: `AES-256-GCM`
- Соль генерируется при регистрации и хранится на сервере
- Шифрование и расшифровка выполняются исключительно на клиенте

### Защита канала

- Канал связи защищён TLS 1.3
- Двухуровневое шифрование: данные шифруются до отправки, затем защищаются TLS

### Восстановление доступа

При смене устройства пользователю достаточно ввести мастер-пароль. Соль и публичный ключ загружаются с сервера, приватный ключ восстанавливается, данные расшифровываются.

## Установка

### Сборка из исходников

```bash
make build-all-platforms
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
# Регистрация
vault register --username=user --password=pass

# Вход
vault login --username=user --password=pass
```

### Команды работы с данными

```bash
# Добавление записи
vault add login --site=example.com --username=user --password=pass

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
./bin/vault.exe register --username=user --password=pass

# Вход
./bin/vault.exe login --username=user --password=pass

# Добавление записи
./bin/vault.exe add login --site=example.com --username=user --password=pass

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
.\bin\vault.exe register --username=user --password=pass

# Вход
.\bin\vault.exe login --username=user --password=pass

# Добавление записи
.\bin\vault.exe add login --site=example.com --username=user --password=pass

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
make test             # Запуск тестов
make test-coverage    # Покрытие кода
make check            # Линтеры и проверки
make build            # Сборка сервера
make build-cli        # Сборка CLI
make build-all-platforms  # Сборка под все платформы
```

## Технологии

- Go 1.26
- gRPC + Protobuf
- PostgreSQL
- Ed25519
- Argon2id
- AES-256-GCM
- TLS 1.3
- Docker
- Cobra

## Лицензия

MIT
```
