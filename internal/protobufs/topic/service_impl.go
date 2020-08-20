package topic

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
	http "github.com/mrdulin/grpc-go-cnode/internal/utils/http"
)

var (
	ErrGetTopicById    = errors.New("topicServiceImpl: get topic by id")
	ErrGetTopicsByPage = errors.New("topicServiceImpl: get topics b")
)

type topicServiceImpl struct {
	HttpClient http.Client
	BaseURL    string
	UnimplementedTopicServiceServer
}

func NewTopicServiceImpl(httpClient http.Client, BaseURL string) *topicServiceImpl {
	return &topicServiceImpl{HttpClient: httpClient, BaseURL: BaseURL}
}

func (t *topicServiceImpl) GetTopicById(ctx context.Context, in *GetTopicByIdRequest) (*GetTopicByIdResponse, error) {
	base, err := url.Parse(t.BaseURL + "/topic/" + in.Id)
	var res GetTopicByIdResponse
	if err != nil {
		fmt.Println("parse url.", err)
		return nil, ErrGetTopicById
	}
	urlValues := url.Values{}
	if in.Accesstoken != "" {
		urlValues.Add("accesstoken", in.Accesstoken)
	}
	urlValues.Add("mdrender", in.Mdrender)
	base.RawQuery = urlValues.Encode()
	err = t.HttpClient.Get(base.String(), &res)
	if err != nil {
		fmt.Println(err)
		return nil, ErrGetTopicById
	}
	return &res, nil
}

func (t *topicServiceImpl) GetTopicsByPage(ctx context.Context, in *GetTopicsByPageRequest) (*GetTopicsByPageResponse, error) {
	base, err := url.Parse(t.BaseURL + "/topics")
	var res GetTopicsByPageResponse
	if err != nil {
		fmt.Println("parse url.", err)
		return nil, ErrGetTopicsByPage
	}
	v, err := query.Values(in)
	if err != nil {
		fmt.Printf("query.Values(params) error. params: %+v, error: %s\n", in, err)
		return nil, ErrGetTopicsByPage
	}
	base.RawQuery = v.Encode()
	err = t.HttpClient.Get(base.String(), &res)
	if err != nil {
		fmt.Println(err)
		return nil, ErrGetTopicsByPage
	}
	return &res, nil
}
