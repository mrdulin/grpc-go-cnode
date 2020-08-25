package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mrdulin/grpc-go-cnode/configs"
	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/topic"
	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"
	"github.com/mrdulin/grpc-go-cnode/internal/utils/http"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	conf *viper.Viper
)

func init() {
	conf = configs.Read()
}

func main() {
	port := conf.GetString(configs.PORT)
	baseurl := conf.GetString(configs.BASE_URL)
	if baseurl == "" {
		log.Fatal("missing api url")
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	httpClient := http.NewClient()
	userServiceImpl := user.NewUserServiceImpl(httpClient, baseurl)
	topicServiceImpl := topic.NewTopicServiceImpl(httpClient, baseurl)
	user.RegisterUserServiceServer(grpcServer, userServiceImpl)
	topic.RegisterTopicServiceServer(grpcServer, topicServiceImpl)
	log.Printf("gRPC server is listening on port: %s\n", port)
	log.Fatal(grpcServer.Serve(lis))
}
