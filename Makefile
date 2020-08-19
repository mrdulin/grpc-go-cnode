
compile-protobuf:
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		./internal/protobufs/user/*.proto

start:
	go run ./cmd/server/main.go
