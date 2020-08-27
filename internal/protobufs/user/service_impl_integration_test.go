package user_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mrdulin/grpc-go-cnode/configs"
	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"
	"github.com/mrdulin/grpc-go-cnode/internal/utils"
	"github.com/mrdulin/grpc-go-cnode/internal/utils/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	serverAddress   = "localhost:3000"
	client          user.UserServiceClient
	clientWithCreds user.UserServiceClient
	conn            *grpc.ClientConn
	connWithCreds   *grpc.ClientConn
	accesstoken     string
)

func setup() {
	fmt.Println("setup")
	conf := configs.Read()
	accesstoken = conf.GetString(configs.ACCESS_TOKEN)
	perRPCCreds := auth.Authentication{Authorization: "Bearer 123"}
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	certFile := path.Join(dir, "../../../assets/server.crt")
	creds, err := credentials.NewClientTLSFromFile(certFile, "localhost")
	if err != nil {
		log.Fatal(err)
	}
	transportCredentials := grpc.WithTransportCredentials(creds)
	conn, err = grpc.Dial(serverAddress, transportCredentials)
	connWithCreds, err = grpc.Dial(serverAddress, transportCredentials, grpc.WithPerRPCCredentials(&perRPCCreds))
	client = user.NewUserServiceClient(conn)
	clientWithCreds = user.NewUserServiceClient(connWithCreds)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	fmt.Println("tearDown")
	conn.Close()
	connWithCreds.Close()
}

func TestUserServiceImpl_GetUserByLoginname_Integration(t *testing.T) {
	utils.MarkedAsIntegrationTest(t)

	t.Run("should get user detail by loginname", func(t *testing.T) {
		args := user.GetUserByLoginnameRequest{Loginname: "mrdulin"}

		res, err := client.GetUserByLoginname(context.Background(), &args)
		if err != nil {
			t.Fatal(err)
		}
		if res.Data == nil {
			t.Errorf("expect get user detail, got: %+v", res)
		}
	})

	t.Run("should return error if the loginname does not exist", func(t *testing.T) {
		args := user.GetUserByLoginnameRequest{Loginname: "1"}

		res, err := client.GetUserByLoginname(context.Background(), &args)
		if res != nil {
			t.Errorf("res should be nil, got: %+v", res)
		}
		if !strings.Contains(err.Error(), "userServiceImpl: Get user by login name") {
			t.Errorf("should get \"userServiceImpl: Get user by login name\" error message, got: %s", err.Error())
		}
	})

	t.Run("should return error if the loginname is empty string", func(t *testing.T) {
		args := user.GetUserByLoginnameRequest{Loginname: ""}

		res, err := client.GetUserByLoginname(context.Background(), &args)
		if res != nil {
			t.Errorf("res should be nil, got: %+v", res)
		}
		s, _ := status.FromError(err)
		if s.Code() != codes.InvalidArgument {
			t.Errorf("should get invalid argument code")
		}
		if s.Message() != "invalid GetUserByLoginnameRequest.Loginname: value length must be at least 1 runes" {
			t.Errorf("should get invalid argument error message")
		}

	})

}

func TestUserServiceImpl_ValidateAccessToken_Integration(t *testing.T) {
	utils.MarkedAsIntegrationTest(t)

	t.Run("should get unauthenticated error for", func(t *testing.T) {
		args := user.ValidateAccessTokenRequest{Accesstoken: accesstoken}

		res, err := client.ValidateAccessToken(context.Background(), &args)
		if res != nil {
			t.Fatalf("expected res is nil, got: %+v", res)
		}
		s, _ := status.FromError(err)
		if s.Code() != codes.Unauthenticated {
			t.Fatalf("expected unanthenticated error code, got: %+v", s.Code())
		}
		if s.Message() != "invalid token" {
			t.Fatalf("expected got invalid token message, got: %s", s.Message())
		}
	})

	t.Run("should validate accesstoken correctly", func(t *testing.T) {
		args := user.ValidateAccessTokenRequest{Accesstoken: accesstoken}
		res, err := clientWithCreds.ValidateAccessToken(context.Background(), &args)
		if err != nil {
			t.Fatal(err)
		}
		if res != nil {
			if !res.GetSuccess() {
				t.Fatalf("expected validate accesstoken success")
			}
			if res.GetLoginname() != "mrdulin" {
				t.Fatalf("expected loginname is mrdulin, got: %s", res.GetLoginname())
			}
		}
	})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
