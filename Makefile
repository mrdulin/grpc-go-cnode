
INTEGRATION_TEST_FUNC_NAME_SUFFIX=Integration

compile-protobuf:
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		./internal/protobufs/**/*.proto

start:
	go run ./cmd/server/main.go

test-integration:
	go test -run ${INTEGRATION_TEST_FUNC_NAME_SUFFIX} -v ./...

test-unit:
	go test -v --short ./...
