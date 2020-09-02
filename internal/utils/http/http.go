package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type ResponseData struct {
	Data json.RawMessage `json:"data,omitempty"`
}
type ResponseStatus struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_msg,omitempty"`
}

// Response API response struct
type Response struct {
	ResponseStatus
	ResponseData
}

type ResponseMap map[string]interface{}

type Decoder interface {
	Decode(body io.ReadCloser, res interface{}) error
}
type Unmarshaler interface {
	Unmarshal(byte interface{}, data proto.Message) error
}
type APIErrorHandler interface {
	HandleAPIError(res interface{}) error
}

type Client interface {
	Get(ctx context.Context, url string, data proto.Message) error
	Post(ctx context.Context, url string, body interface{}, data proto.Message) error
}

type client struct{}

func NewClient() *client {
	return new(client)
}

func createHttpClient() *http.Client {
	return &http.Client{
		Transport: &ochttp.Transport{
			Propagation: &propagation.HTTPFormat{},
		},
	}
}

//Get send GET HTTP request
func (h *client) Get(ctx context.Context, url string, data proto.Message) error {
	httpClient := createHttpClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest(\"GET\", url, nil)")
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do(req)")
	}
	defer resp.Body.Close()
	var res Response
	if err = h.Decode(resp.Body, &res); err != nil {
		return err
	}
	if err = h.HandleAPIError(res); err != nil {
		return err
	}
	if err = h.Unmarshal(res, data); err != nil {
		return err
	}
	return nil
}

//Post send POST HTTP request
func (h *client) Post(ctx context.Context, url string, body interface{}, data proto.Message) error {
	var res ResponseMap
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return errors.Wrapf(err, "json.Marshal(body). body: %+v", body)
	}
	httpClient := createHttpClient()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.Wrap(err, "http.NewRequest(\"POST\", url, bytes.NewBuffer(jsonValue))")
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "httpClient.Do(req) req: %+v", req)
	}
	defer resp.Body.Close()
	if err = h.Decode(resp.Body, &res); err != nil {
		return err
	}
	if err = h.HandleAPIError(res); err != nil {
		return err
	}
	if err = h.Unmarshal(res, data); err != nil {
		return err
	}
	return nil
}

func (h *client) Decode(body io.ReadCloser, res interface{}) error {
	err := json.NewDecoder(body).Decode(res)
	if err != nil {
		return errors.Wrapf(err, "json.NewDecoder(resp.Body).Decode(&res)")
	}
	return nil
}

func (h *client) Unmarshal(res interface{}, data proto.Message) error {
	bs, err := json.Marshal(res)
	if err != nil {
		return errors.Wrapf(err, "json.Marshal(r). res: %+v", res)
	}
	um := jsonpb.Unmarshaler{}
	err = um.Unmarshal(strings.NewReader(string(bs)), data)
	if err != nil {
		return errors.Wrapf(err, "json.Unmarshal. data: %s", string(bs))
	}
	return nil
}

func (h *client) HandleAPIError(res interface{}) error {
	var (
		success      bool
		errorMessage string
	)
	switch v := res.(type) {
	case Response:
		success = v.Success
		errorMessage = v.ErrorMessage
	case ResponseMap:
		success = v["success"].(bool)
		if v["error_msg"] != nil {
			errorMessage = v["error_msg"].(string)
		}
	}
	if !success {
		return fmt.Errorf("API error: %s", errorMessage)
	}
	return nil
}
