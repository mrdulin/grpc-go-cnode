# build
FROM golang:1.14 as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build

# release
#FROM scratch
FROM alpine:3.7
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/bin/app /app/
EXPOSE 3000
ENTRYPOINT ["/app/app"]
