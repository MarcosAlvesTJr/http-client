package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HttpClient struct {
	BaseURL string
	client  http.Client
	headers map[string]string
}

type Response struct {
	httpResponse *http.Response
}

func (r *Response) DecodeTo(v any) error {
	responseDecoder := json.NewDecoder(r.httpResponse.Body)

	defer closeBody(r.httpResponse.Body)

	err := responseDecoder.Decode(&v)

	if err != nil {
		return err
	}

	return nil
}

func NewHttpClient(baseUrl string) *HttpClient {
	defaultHeaders := make(map[string]string)

	defaultHeaders["Content-Type"] = "application/json"
	defaultHeaders["Accept"] = "application/json"

	return &HttpClient{
		BaseURL: formatBaseURL(baseUrl),
		client:  http.Client{},
		headers: defaultHeaders,
	}
}

func (c *HttpClient) AddHeader(key, value string) {
	c.headers[key] = value
}

func (c *HttpClient) Get(url string) (*Response, error) {
	req, err := http.NewRequest("GET", formatUrl(url, c.BaseURL), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(requestInterceptor(req, c.headers))

	if err != nil {
		return nil, err
	}

	return &Response{httpResponse: res}, nil
}

func (c *HttpClient) Post(url string, body any) (*Response, error) {
	jsonBytes, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", formatUrl(url, c.BaseURL), bytes.NewBuffer(jsonBytes))

	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(requestInterceptor(req, c.headers))

	if err != nil {
		return nil, err
	}

	return &Response{httpResponse: res}, nil
}

func (c *HttpClient) Delete(url string) (*Response, error) {
	req, err := http.NewRequest("DELETE", formatUrl(url, c.BaseURL), nil)

	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(requestInterceptor(req, c.headers))

	if err != nil {
		return nil, err
	}

	return &Response{httpResponse: res}, nil
}

func (c *HttpClient) Put(url string, body any) (*Response, error) {
	jsonBytes, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", formatUrl(url, c.BaseURL), bytes.NewBuffer(jsonBytes))

	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(requestInterceptor(req, c.headers))

	if err != nil {
		return nil, err
	}

	return &Response{httpResponse: res}, nil
}

func requestInterceptor(request *http.Request, headers map[string]string) *http.Request {
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	return request
}

func formatBaseURL(url string) string {
	if strings.HasSuffix(url, "/") {
		return url
	}

	return fmt.Sprintf("%s/", url)
}

func formatUrl(url, baseUrl string) string {
	return fmt.Sprintf("%s%s", baseUrl, strings.TrimPrefix(url, "/"))
}

func closeBody(body io.ReadCloser) {
	err := body.Close()

	if err != nil {
		panic(err)
	}
}
