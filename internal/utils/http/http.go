package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/golang/protobuf/jsonpb"

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

type Client interface {
	Get(url string, data proto.Message) error
	Post(url string, body interface{}, data proto.Message) error
	HandleAPIError(res interface{}) error
	Decode(body io.ReadCloser, res interface{}) error
	Unmarshal(byte interface{}, data proto.Message) error
}

type client struct{}

func NewClient() *client {
	return &client{}
}

//Get send GET HTTP request
func (h *client) Get(url string, data proto.Message) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "http.Get(url)")
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
func (h *client) Post(url string, body interface{}, data proto.Message) error {
	var res ResponseMap
	jsonValue, err := json.Marshal(body)
	if err != nil {
		return errors.Wrapf(err, "json.Marshal(body). body: %+v", body)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return errors.Wrapf(err, "http.Post(url, \"application/json\", bytes.NewBuffer(jsonValue)). jsonValue: %+v", jsonValue)
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
	var (
		bs  []byte
		err error
	)
	switch v := res.(type) {
	case Response:
		bs = v.Data
	case ResponseMap:
		var r interface{}
		if v["data"] != nil {
			r = v["data"]
		}
		r = v
		bs, err = json.Marshal(r)
		if err != nil {
			return errors.Wrapf(err, "json.Marshal(r). v: %+v", r)
		}
	}
	um := jsonpb.Unmarshaler{}
	err = um.Unmarshal(strings.NewReader(string(bs)), data)
	//err = json.Unmarshal(bs, &data)
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
