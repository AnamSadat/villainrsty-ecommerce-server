package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"villainrsty-ecommerce-server/internal/core/shared/errors"
)

type (
	HTTPClient struct {
		client  *http.Client
		baseURL string
	}

	Request struct {
		Method  string
		Path    string
		Headers map[string]string
		Body    interface{}
	}

	Response struct {
		StatusCode int
		Body       []byte
		Headers    http.Header
	}
)

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

func (c *HTTPClient) Do(req Request) (*Response, error) {
	url := c.baseURL + req.Path

	var body io.Reader
	if req.Body == nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "Failed to marshal request body", err)
		}

		body = bytes.NewReader(bodyBytes)
	}

	httpReq, err := http.NewRequest(req.Method, url, body)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "Failed to create request", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "Failed to do reques", err)
	}

	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "Failed to read response body", err)
	}

	return &Response{
		StatusCode: httpResp.StatusCode,
		Body:       respBody,
		Headers:    httpResp.Header,
	}, nil
}

func (c *HTTPClient) Get(path string, headers map[string]string) (*Response, error) {
	return c.Do(Request{
		Method:  http.MethodGet,
		Path:    path,
		Headers: headers,
	})
}

func (c *HTTPClient) Post(path string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Do(Request{
		Method:  http.MethodPost,
		Path:    path,
		Body:    body,
		Headers: headers,
	})
}

func (c *HTTPClient) Put(path string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Do(Request{
		Method:  http.MethodPut,
		Path:    path,
		Body:    body,
		Headers: headers,
	})
}

func (c *HTTPClient) Patch(path string, body interface{}, headers map[string]string) (*Response, error) {
	return c.Do(Request{
		Method:  http.MethodPatch,
		Path:    path,
		Body:    body,
		Headers: headers,
	})
}

func (c *HTTPClient) Delete(path string, headers map[string]string) (*Response, error) {
	return c.Do(Request{
		Method:  http.MethodDelete,
		Path:    path,
		Headers: headers,
	})
}
