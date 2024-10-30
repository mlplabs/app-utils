package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const Token string = "TOKEN"

type RequestParams struct {
	ProxyToken     bool
	RequestBody    interface{}
	ResponseBody   interface{}
	RequestHandler func(request *http.Request) *http.Request
}

type Client struct {
	client           *http.Client
	clientName       string // какой сервис выполняет запрос
	ownerServiceName string // какому сервису запрос
	baseURL          string
	body             []byte
}

func NewClient(clientName string, ownerServiceName string, baseURL string) *Client {
	return &Client{
		client:           &http.Client{},
		clientName:       clientName,
		ownerServiceName: ownerServiceName,
		baseURL:          baseURL,
	}
}

func (c *Client) GetBaseURL() string {
	return c.baseURL
}
func (c *Client) Get(ctx context.Context, url string, params *RequestParams) (*http.Response, error) {
	return c.get(ctx, http.MethodGet, url, params)
}
func (c *Client) Delete(ctx context.Context, url string, params *RequestParams) (*http.Response, error) {
	return c.get(ctx, http.MethodDelete, url, params)
}
func (c *Client) Post(ctx context.Context, url string, params *RequestParams) (*http.Response, error) {
	return c.post(ctx, http.MethodPost, url, params)
}
func (c *Client) Put(ctx context.Context, url string, params *RequestParams) (*http.Response, error) {
	return c.post(ctx, http.MethodPut, url, params)
}
func (c *Client) get(ctx context.Context, method string, url string, params *RequestParams) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.GetBaseURL(), url), nil)
	if err != nil {
		return nil, err
	}
	if params != nil {
		if params.ProxyToken {
			token, ok := c.GetToken(ctx)
			if !ok {
				return nil, fmt.Errorf("no token in context")
			}
			req.Header.Add("Authorization", "Bearer "+token)
		}

		if params.RequestHandler != nil {
			req = params.RequestHandler(req)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusBadRequest {
		if params != nil && params.ResponseBody != nil {
			err = c.ReadBody(resp)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(c.body, &params.ResponseBody)
			if err != nil {
				return nil, err
			}
		}
		return resp, err
	} else {
		return resp, fmt.Errorf("server error, code: %d service: %s", resp.StatusCode, c.ownerServiceName)
	}
}
func (c *Client) post(ctx context.Context, method string, url string, params *RequestParams) (*http.Response, error) {
	var bodyBuffer *bytes.Buffer
	if params != nil && params.RequestBody != nil {
		body, err := json.Marshal(params.RequestBody)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.GetBaseURL(), url), bodyBuffer)
	if err != nil {
		return nil, err
	}

	if params != nil {
		if params.ProxyToken {
			token, ok := c.GetToken(ctx)
			if !ok {
				return nil, fmt.Errorf("no token in context")
			}
			req.Header.Add("Authorization", token)
		}

		if params.RequestHandler != nil {
			req = params.RequestHandler(req)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusBadRequest {
		if params != nil && params.ResponseBody != nil {
			err = c.ReadBody(resp)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(c.body, &params.ResponseBody)
			if err != nil {
				return nil, err
			}
		}
		return resp, err
	} else {
		return resp, fmt.Errorf("server error, code: %d service: %s", resp.StatusCode, c.ownerServiceName)
	}

}
func (c *Client) ReadBody(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.body = body
	resp.Body.Close()
	return nil
}
func (c *Client) GetToken(ctx context.Context) (string, bool) {
	val := ctx.Value(Token)
	value, ok := val.(string)
	return value, ok
}
