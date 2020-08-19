package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/topic"

	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"
	"google.golang.org/grpc"
)

const (
	accesstoken string = "be60f8d0-149c-4905-be4a-7f07d4788d88"
)

func main() {
	serverAddress := "localhost:3000"

	conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())

	if e != nil {
		panic(e)
	}
	defer conn.Close()

	topicServiceClient := topic.NewTopicServiceClient(conn)
	testGetTopicById(topicServiceClient)

	//userServiceClient := user.NewUserServiceClient(conn)
	//testValidateAccessToken(userServiceClient)
}

func testGetUserByLoginname(client user.UserServiceClient) {
	args := user.GetUserByLoginnameRequest{Loginname: "mrdulin"}

	res, err := client.GetUserByLoginname(context.Background(), &args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GetUserByLoginname: %+v", res)
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
