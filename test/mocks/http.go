package mocks

import (
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/mock"
)

type MockedHttp struct {
	mock.Mock
}

func (mh *MockedHttp) Get(url string, data proto.Message) error {
	args := mh.Called(url, data)
	return args.Error(0)
}

func (mh *MockedHttp) Post(url string, body interface{}, data proto.Message) error {
	args := mh.Called(url, body, data)
	return args.Error(0)
}
