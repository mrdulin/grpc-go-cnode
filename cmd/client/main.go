package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/credentials"

	"github.com/mrdulin/grpc-go-cnode/configs"

	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/topic"

	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"
	"google.golang.org/grpc"
)

var (
	accesstoken   string
	serverAddress string
)

func init() {
	conf := configs.Read()
	accesstoken = conf.GetString(configs.ACCESS_TOKEN)
	serverAddress = fmt.Sprintf("localhost:%s", conf.GetString(configs.PORT))
}

func main() {
	creds, err := credentials.NewClientTLSFromFile("./assets/server.crt", "server.grpc.io")
	if err != nil {
		log.Fatal(err)
	}
	conn, e := grpc.Dial(serverAddress,
		//grpc.WithInsecure()
		grpc.WithTransportCredentials(creds),
	)

	if e != nil {
		panic(e)
	}
	defer conn.Close()

	topicServiceClient := topic.NewTopicServiceClient(conn)
	testGetTopicById(topicServiceClient)
	//testGetTopicsByPage(topicServiceClient)
}

func testValidateAccessToken(client user.UserServiceClient) {
	args := user.ValidateAccessTokenRequest{Accesstoken: accesstoken}
	res, err := client.ValidateAccessToken(context.Background(), &args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ValidateAccessToken: %+v", res)
}

func testGetTopicById(client topic.TopicServiceClient) {
	args := topic.GetTopicByIdRequest{Id: "5433d5e4e737cbe96dcef312"}
	res, err := client.GetTopicById(context.Background(), &args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GetTopicById: %+v", res)
}

func testGetTopicsByPage(client topic.TopicServiceClient) {
	args := topic.GetTopicsByPageRequest{Page: 1, Limit: 1, Mdrender: "true"}
	res, err := client.GetTopicsByPage(context.Background(), &args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GetTopicsByPage: %+v", res)
}
