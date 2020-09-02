package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
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
	// Register stats and trace exporters to export
	// the collected data.
	view.RegisterExporter(&exporter.PrintExporter{})
	view.SetReportingPeriod(time.Second)

	// Register the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewClientTLSFromFile("./assets/server.crt", "localhost")
	if err != nil {
		log.Fatal(err)
	}
	conn, e := grpc.Dial(serverAddress,
		//grpc.WithInsecure()
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithTransportCredentials(creds),
	)

	if e != nil {
		panic(e)
	}
	defer conn.Close()

	topicServiceClient := topic.NewTopicServiceClient(conn)
	// test
	for {
		testGetTopicById(topicServiceClient)
		time.Sleep(5 * time.Second)
	}
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
