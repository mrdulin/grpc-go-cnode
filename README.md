# grpc-go-cnode

A Go gRPC server for [CNode](https://cnodejs.org) community built on the top of RESTful API.


After modifying the `.proto` files, need to re-compile protocol buffers:

```bash
make compile-protobuf
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
