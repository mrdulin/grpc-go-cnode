package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"
	"google.golang.org/grpc"
)

func main() {
	serverAddress := "localhost:3000"

	conn, e := grpc.Dial(serverAddress, grpc.WithInsecure())

	if e != nil {
		panic(e)
	}
	defer conn.Close()

	client := user.NewUserServiceClient(conn)

	args := user.GetUserByLoginnameRequest{Loginname: "mrdulin"}

	res, err := client.GetUserByLoginname(context.Background(), &args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GetUserByLoginnameResponse: %+v", res)
}
