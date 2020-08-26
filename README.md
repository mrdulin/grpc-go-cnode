# grpc-go-cnode

A Go gRPC server for [CNode](https://cnodejs.org) community built on the top of RESTful API.


After modifying the `.proto` files, need to re-compile protocol buffers.
It will generate the service interfaces, models for the server side and service interfaces for client side.
Then, you can implement the interfaces of the services.

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

Start the gRPC server:
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
