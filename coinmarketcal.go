package coinmarketcal

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

const DefaultURL = "https://developers.coinmarketcal.com/v1"
const CustomKeyHeader = "x-api-key"

type Client struct {
	ApiUrl string
	ApiKey string
	c      *http.Client
}

func NewClient(url, key string) *Client {
	return &Client{
		ApiUrl: url,
		ApiKey: key,
		c:      http.DefaultClient,
	}
}

func (c *Client) Url(endpoint string) string {
	return fmt.Sprintf("%s%s", c.ApiUrl, endpoint)
}

func (c *Client) GetCall(ctx context.Context, endpoint string, opt, response any) (result *http.Response, err error) {
	var path string
	if opt != nil {
		v, err := query.Values(opt)
		if err != nil {
			return result, err
		}
		path = fmt.Sprintf("%s?%s", c.Url(endpoint), v.Encode())
	} else {
		path = c.Url(endpoint)
	}
	// Parse the base URL
	parsedURL, err := url.Parse(path)
	if err != nil {
		return result, err
	}
	return c.doCall(ctx, NewRequest(parsedURL.String(), nil), response)
}

func (c *Client) doCall(ctx context.Context, req *Request, response any) (*http.Response, error) {
	if reflect.TypeOf(response).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("response struct is not a pointer")
	}
	httpRequest, err := req.NewHttpRequest(ctx, c.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("api call %v() on %v: %v", req.Method, req.Endpoint, err.Error())
	}
	httpResponse, err := c.c.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("api call %v() on %v: %v", req.Method, httpRequest.URL.String(), err.Error())
	}

	var reader io.ReadCloser
	switch httpResponse.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(httpResponse.Body)
		if err != nil {
			return nil, fmt.Errorf(
				"call %v() on %v status code: %v. could not create gzip reader",
				req.Method,
				httpRequest.URL.String(),
				httpResponse.StatusCode)
		}
	case "deflate":
		reader, err = zlib.NewReader(httpResponse.Body)
		if err != nil {
			return nil, fmt.Errorf(
				"call %v() on %v status code: %v. could not create deflate reader",
				req.Method,
				httpRequest.URL.String(),
				httpResponse.StatusCode)
		}
	default:
		reader = httpResponse.Body
	}
	defer reader.Close()

	bodyBytes, err := io.ReadAll(reader)
	// parsing error
	if err != nil {
		return nil, fmt.Errorf(
			"call %v() on %v status code: %v. could not decode body to response: %v",
			req.Method,
			httpRequest.URL.String(),
			httpResponse.StatusCode,
			err.Error())
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf(
			"call %v() on %v status code: %v",
			req.Method,
			httpRequest.URL.String(),
			httpResponse.StatusCode)
	}

	// 2. Create a gzip reader, wrapping the bytes.Reader.
	err = json.Unmarshal(bodyBytes, response)
	if err != nil {
		return nil, fmt.Errorf(
			"call %v() on %v status code: %v. could not decode body to response model: %v",
			req.Method,
			httpRequest.URL.String(),
			httpResponse.StatusCode,
			err.Error())
	}
	return httpResponse, nil
}

type Request struct {
	Body     any
	Endpoint string
	Method   string
}

type DefaultBodyStatus struct {
	ErrorCode    int64       `json:"error_code"`
	ErrorMessage json.Number `json:"error_message"`
}

type DefaultArrayResponse[T any] struct {
	Data []T `json:"data"`
}

type DefaultBody[T any] struct {
	Body     T                  `json:"body,omitempty"`
	Metadata *Metadata          `json:"metadata,omitempty"`
	Status   *DefaultBodyStatus `json:"status,omitempty"`
}

type Metadata struct {
	Max        int `json:"max"`
	Page       int `json:"page"`
	PageCount  int `json:"page_count"`
	TotalCount int `json:"total_count"`
}

func NewRequest(endpoint string, body any, methods ...string) *Request {
	method := http.MethodGet
	if len(methods) != 0 {
		method = methods[0]
	}
	return &Request{
		Body:     body,
		Endpoint: endpoint,
		Method:   method,
	}
}

func (r *Request) NewHttpRequest(ctx context.Context, token string) (*http.Request, error) {
	var bodyReader io.Reader
	if r.Method != http.MethodGet {
		byteBody, err := json.Marshal(r.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(byteBody)
	}
	request, err := http.NewRequest(r.Method, r.Endpoint, bodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Accept-Encoding", "deflate, gzip")
	request.Header.Add("Accept", "application/json")
	request.Header.Set(CustomKeyHeader, token)
	//Custom headers
	request = request.WithContext(ctx)

	return request, nil
}
