package clover

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OAuthClient handles Clover OAuth-specific endpoints.
// It uses separate base URLs from the REST API client.
type OAuthClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewOAuthClient creates an OAuth client for the given environment.
func NewOAuthClient(env string) *OAuthClient {
	var baseURL string
	switch env {
	case "production", "prod":
		baseURL = "https://api.clover.com"
	case "sandbox", "dev", "development":
		baseURL = "https://apisandbox.dev.clover.com"
	default:
		baseURL = "https://apisandbox.dev.clover.com"
	}

	return &OAuthClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// ExchangeCode exchanges an authorization code for access and refresh tokens.
func (c *OAuthClient) ExchangeCode(ctx context.Context, clientID, clientSecret, code string) (*TokenResponse, error) {
	payload := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}
	return c.postToken(ctx, "/oauth/v2/token", payload)
}

// RefreshTokens obtains a new access and refresh token pair using a refresh token.
// Note: The refresh token is single-use and will be invalidated after this call.
func (c *OAuthClient) RefreshTokens(ctx context.Context, clientID, refreshToken string) (*TokenResponse, error) {
	payload := map[string]string{
		"client_id":     clientID,
		"refresh_token": refreshToken,
	}
	return c.postToken(ctx, "/oauth/v2/refresh", payload)
}

func (c *OAuthClient) postToken(ctx context.Context, path string, payload map[string]string) (*TokenResponse, error) {
	url := c.baseURL + path

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("oauth error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, fmt.Errorf("unmarshal token response: %w", err)
	}

	return &tokenResp, nil
}

// AuthorizeURL returns the Clover OAuth authorization URL for the given environment.
// This is used by the frontend to redirect the merchant to Clover.
func AuthorizeURL(env, clientID, redirectURI string) string {
	var baseURL string
	switch env {
	case "production", "prod":
		baseURL = "https://www.clover.com"
	case "sandbox", "dev", "development":
		baseURL = "https://sandbox.dev.clover.com"
	default:
		baseURL = "https://sandbox.dev.clover.com"
	}

	return fmt.Sprintf("%s/oauth/v2/authorize?client_id=%s&redirect_uri=%s",
		baseURL, clientID, redirectURI)
}
