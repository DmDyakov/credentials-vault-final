FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY server/go.mod server/go.sum ./server/
COPY client-cli/go.mod client-cli/go.sum ./client-cli/
COPY gen/go.mod gen/go.sum ./gen/
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY tools/staticlint/go.mod tools/staticlint/go.sum ./tools/staticlint/

RUN go mod download

COPY server/ ./server/
COPY gen/ ./gen/
COPY pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -C server -o /server ./cmd/server

FROM alpine:3.24

RUN apk add --no-cache ca-certificates

COPY --from=builder /server /server

EXPOSE 9090

ENTRYPOINT ["/server"]
