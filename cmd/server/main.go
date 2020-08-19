package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mrdulin/grpc-go-cnode/internal/utils/http"

	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"
	"google.golang.org/grpc"
)

const (
	port    string = "3000"
	baseurl string = "https://cnodejs.org/api/v1"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	httpClient := http.NewClient()
	userServiceImpl := user.NewUserServiceImpl(httpClient, baseurl)
	user.RegisterUserServiceServer(grpcServer, userServiceImpl)
	log.Fatal(grpcServer.Serve(lis))
}
