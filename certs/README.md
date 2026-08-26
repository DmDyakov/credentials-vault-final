# TLS Сертификаты

## Генерация для dev

### Вариант 1: openssl (без дополнительных утилит)

#### Linux / macOS

```bash
cd certs

# Приватный ключ
openssl genrsa -out server.key 2048

# Публичный самоподписанный сертификат с SAN
openssl req -x509 -newkey rsa:2048 -key server.key -out server.crt -days 365 \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

cd ..
```

#### Windows (Git Bash)

```bash
cd certs

openssl genrsa -out server.key 2048

openssl req -x509 -newkey rsa:2048 -key server.key -out server.crt -days 365 \
  -subj "//CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

cd ..
```

#### Windows (PowerShell)

```powershell
cd certs

openssl genrsa -out server.key 2048

openssl req -x509 -newkey rsa:2048 -key server.key -out server.crt -days 365 `
  -subj "/CN=localhost" `
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

cd ..
```

### Вариант 2: mkcert (доверенные сертификаты)

mkcert создаёт локальный CA и сертификаты, которым доверяет система без предупреждений.

#### Установка

```bash
# macOS
brew install mkcert

# Windows (winget)
winget install FiloSottile.mkcert

# Windows (choco)
choco install mkcert

# Ubuntu/Debian
sudo apt install mkcert

# Fedora
sudo dnf install mkcert

# Arch Linux
sudo pacman -S mkcert
```

#### Генерация

```bash
# Установить локальный CA (один раз)
mkcert -install

# Сгенерировать сертификат
cd certs
mkcert localhost 127.0.0.1

# mkcert создаст файлы:
# localhost+1.pem — сертификат
# localhost+1-key.pem — ключ
```

## Файлы

- `server.crt` — публичный сертификат (не коммитить)
- `server.key` — приватный ключ (не коммитить)
- Или файлы от mkcert: `localhost+1.pem`, `localhost+1-key.pem`

## Использование

- **Сервер**: `TLS_CERT_FILE=certs/server.crt`, `TLS_KEY_FILE=certs/server.key`
- **Клиент**: `CAFile=certs/server.crt`

## Для CI

Сертификаты генерируются автоматически в CI-пайплайне:

```yaml
- name: Generate TLS certs
  run: |
    cd certs
    openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 1 -nodes \
      -subj "/CN=localhost" \
      -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
```

## Для staging и prod

Использовать сертификаты от Let's Encrypt или коммерческого CA.
Клиент не указывает `CAFile` — используются системные корневые CA.
