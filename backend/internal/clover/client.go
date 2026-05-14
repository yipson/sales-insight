package clover

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a generic HTTP client for the Clover API.
type Client struct {
	httpClient   *http.Client
	baseURL      string
	rateLimiter  *RateLimiter
}

// NewClient creates a Clover API client for the given environment.
// env can be "sandbox" or "production".
func NewClient(env string) *Client {
	var baseURL string
	switch env {
	case "production", "prod":
		baseURL = "https://api.clover.com/v3"
	case "sandbox", "dev", "development":
		baseURL = "https://sandbox.dev.clover.com/v3"
	default:
		baseURL = "https://sandbox.dev.clover.com/v3"
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:     baseURL,
		rateLimiter: NewRateLimiter(),
	}
}

// WithTimeout returns a new client with the specified timeout.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: timeout},
		baseURL:     c.baseURL,
		rateLimiter: c.rateLimiter,
	}
}

// Get performs a GET request to the Clover API.
func (c *Client) Get(ctx context.Context, path, token string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, path, token, nil)
}

// Post performs a POST request to the Clover API.
func (c *Client) Post(ctx context.Context, path, token string, body []byte) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPost, path, token, body)
}

// Put performs a PUT request to the Clover API.
func (c *Client) Put(ctx context.Context, path, token string, body []byte) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPut, path, token, body)
}

// Delete performs a DELETE request to the Clover API.
func (c *Client) Delete(ctx context.Context, path, token string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodDelete, path, token, nil)
}

func (c *Client) doRequest(ctx context.Context, method, path, token string, body []byte) ([]byte, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.rateLimiter.Do(c.httpClient, req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("clover API error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
