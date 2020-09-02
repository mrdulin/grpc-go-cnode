package mocks

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/mock"
)

type MockedHttp struct {
	mock.Mock
}

func (mh *MockedHttp) Get(ctx context.Context, url string, data proto.Message) error {
	args := mh.Called(ctx, url, data)
	return args.Error(0)
}

func (mh *MockedHttp) Post(ctx context.Context, url string, body interface{}, data proto.Message) error {
	args := mh.Called(ctx, url, body, data)
	return args.Error(0)
}
