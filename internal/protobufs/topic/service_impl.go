package topic

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	http "github.com/mrdulin/grpc-go-cnode/internal/utils/http"
)

var (
	ErrGetTopicById = errors.New("topicServiceImpl: get topic by id")
)

type topicServiceImpl struct {
	HttpClient http.Client
	BaseURL    string
	UnimplementedTopicServiceServer
}

func NewTopicServiceImpl(httpClient http.Client, BaseURL string) *topicServiceImpl {
	return &topicServiceImpl{HttpClient: httpClient, BaseURL: BaseURL}
}

func (t *topicServiceImpl) GetTopicById(ctx context.Context, in *GetTopicByIdRequest) (*TopicDetail, error) {
	base, err := url.Parse(t.BaseURL + "/topic/" + in.Id)
	var res TopicDetail
	if err != nil {
		fmt.Println("Get topic by id error: parse url.", err)
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
