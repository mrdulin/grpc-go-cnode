package user

import (
	"context"
	"errors"
	"fmt"

	http "github.com/mrdulin/grpc-go-cnode/internal/utils/http"
)

var (
	ErrGetUserByLoginname  = errors.New("get user by login name")
	ErrValidateAccessToken = errors.New("Validate accessToken")
)

type userServiceImpl struct {
	HttpClient http.Client
	BaseURL    string
	UnimplementedUserServiceServer
}

func NewUserServiceImpl(httpClient http.Client, baseurl string) *userServiceImpl {
	return &userServiceImpl{HttpClient: httpClient, BaseURL: baseurl}
}

func (svc *userServiceImpl) GetUserByLoginname(ctx context.Context, in *GetUserByLoginnameRequest) (*User, error) {
	endpoint := svc.BaseURL + "/user/" + in.Loginname
	var res User
	err := svc.HttpClient.Get(endpoint, &res)
	if err != nil {
		fmt.Println(err)
		return nil, ErrGetUserByLoginname
	}
	return &res, nil
}
