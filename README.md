# grpc-go-cnode

[![Go Report Card](https://goreportcard.com/badge/github.com/mrdulin/grpc-go-cnode)](https://goreportcard.com/report/github.com/mrdulin/grpc-go-cnode)

A Go gRPC server for [CNode](https://cnodejs.org) community built on the top of RESTful API.


After modifying the `.proto` files, need to re-compile protocol buffers.
It will generate the service interfaces, models for the server side and service interfaces for client side.
Then, you can implement the interfaces of the services.

Features:

* HTTPs server and gRPC server share same listening address and port.
* Print access logs in unary call interceptor
* gRPC Health check for all services based on [GRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
* Per RPC call authentication, check [auth.go](./internal/utils/auth/auth.go)
* TLS connection with self-signed credentials
* Support constraint rules and validators for Protocol buffer, check [here](./internal/protobufs/user/service.proto)

Compile protocol buffers:

```bash
make compile-protobuf
```

Environment variables in `configs/config.yaml`:
```
BASE_URL: https://cnodejs.org/api/v1
PORT: 3000
ACCESS_TOKEN: <YOUR_ACCESS_TOKEN>
GRPC_GO_LOG_SEVERITY_LEVEL: info
GRPC_GO_LOG_VERBOSITY_LEVEL: 1
```

Start the HTTPs server and gRPC server:
```bash
make start
```

Run integration testings:

1. Run `make start` to start the server
2. Run `make test-integration`

Run unit testings:

```bash
make test-unit
```

More info, see [Makefile](./Makefile)
