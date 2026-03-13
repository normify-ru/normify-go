package normify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     apiKey,
		baseURL:    "https://normify.ru",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func (c *Client) Process(ctx context.Context, req *ProcessRequest) (*ProcessResponse, error) {
	url := c.baseURL + "/api/v1/auth/process"

	if len(req.Data) > 10 {
		return nil, fmt.Errorf("req data is too long")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа для последующего анализа
	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       respBody.Bytes(),
		}
	}

	var processResp ProcessResponse
	if err := json.Unmarshal(respBody.Bytes(), &processResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &processResp, nil
}
