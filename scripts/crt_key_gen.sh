#!/bin/sh

echo "generating certificates and private key for server and client..."
openssl genrsa -out ./assets/server.key 2048

openssl req -new -x509 -days 3650 \
  -subj "/C=GB/L=China/O=grpo-server/CN=localhost" \
  -key ./assets/server.key -out ./assets/server.crt

openssl genrsa -out ./assets/client.key 2048

openssl req -new -x509 -days 3650 \
  -subj "/C=GB/L=China/O=grpo-client/CN=localhost" \
  -key ./assets/client.key -out ./assets/client.crt
