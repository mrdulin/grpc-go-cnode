package user_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/mrdulin/grpc-go-cnode/internal/protobufs/user"

	"google.golang.org/grpc"

	"github.com/mrdulin/grpc-go-cnode/internal/utils"
)

var (
	serverAddress = "localhost:3000"
	client        user.UserServiceClient
	conn          *grpc.ClientConn
	err           error
)

func setup() {
	fmt.Println("setup")
	conn, err = grpc.Dial(serverAddress, grpc.WithInsecure())
	client = user.NewUserServiceClient(conn)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	fmt.Println("tearDown")
	conn.Close()
}

func TestUserServiceImpl_GetUserByLoginname_Integration(t *testing.T) {
	utils.MarkedAsIntegrationTest(t)

	t.Run("should get user detail by loginname", func(t *testing.T) {
		args := user.GetUserByLoginnameRequest{Loginname: "mrdulin"}

		res, err := client.GetUserByLoginname(context.Background(), &args)
		if err != nil {
			log.Fatal(err)
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
		t.Error(err)
	})

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
