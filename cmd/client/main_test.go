package main_test

import (
	"context"
	"log"
	"os"
	"path"
	"testing"

	"github.com/mrdulin/grpc-go-cnode/internal/utils"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	serverAddress = "localhost:3000"
)

func TestHealthCheckService_Integration(t *testing.T) {
	utils.MarkedAsIntegrationTest(t)
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	certFile := path.Join(dir, "../../assets/server.crt")
	creds, err := credentials.NewClientTLSFromFile(certFile, "localhost")
	if err != nil {
		log.Fatal(err)
	}
	conn, e := grpc.Dial(serverAddress, grpc.WithTransportCredentials(creds))

	if e != nil {
		panic(e)
	}
	defer conn.Close()

	hcClient := healthpb.NewHealthClient(conn)
	req := new(healthpb.HealthCheckRequest)
	res, err := hcClient.Check(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if res.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Errorf("expect serving status, got: %+v", res.GetStatus())
	}
}
