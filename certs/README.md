```markdown
# TLS Сертификаты

## Генерация для dev

### Linux

```bash
cd certs

# Приватный ключ (секрет)
openssl genrsa -out server.key 2048

# Публичный самоподписанный сертификат с SAN
openssl req -x509 -newkey rsa:2048 -key server.key -out server.crt -days 365 \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

cd ..
```

### macOS

```bash
cd certs

# Приватный ключ (секрет)
openssl genrsa -out server.key 2048

# Публичный самоподписанный сертификат с SAN
openssl req -x509 -newkey rsa:2048 -key server.key -out server.crt -days 365 \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

cd ..
```

### Windows (Git Bash)

```bash
cd certs

# Приватный ключ (секрет)
openssl genrsa -out server.key 2048

# Публичный самоподписанный сертификат с SAN (обратите внимание на //CN)
openssl req -x509 -newkey rsa:2048 -key server.key -out server.crt -days 365 \
  -subj "//CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

cd ..
```

### Windows (PowerShell)

```powershell
cd certs

# Приватный ключ (секрет)
openssl genrsa -out server.key 2048

# Публичный самоподписанный сертификат с SAN
openssl req -x509 -newkey rsa:2048 -key server.key -out server.crt -days 365 `
  -subj "/CN=localhost" `
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

cd ..
```

## Файлы

- `server.crt` — публичный сертификат (можно коммитить)
- `server.key` — приватный ключ (НЕ коммитить, добавить в .gitignore)

## Использование

- **Сервер**: `TLS_CERT_FILE=/certs/server.crt`, `TLS_KEY_FILE=/certs/server.key`
- **Клиент**: `CAFile=certs/server.crt` (доверяет сертификату)

## Для staging и prod

Использовать сертификаты от Let's Encrypt или коммерческого CA.
Клиент не указывает `CAFile` — используются системные корневые CA.
```
