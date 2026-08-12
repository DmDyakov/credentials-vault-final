.PHONY: proto
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