package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LambdaToolClient makes authenticated HTTP requests to a lambda tool endpoint.
type LambdaToolClient struct {
	url    string
	apiKey string
}

// NewLambdaToolClient creates a client for the given endpoint URL and API key.
// apiKey may be empty for local development (no API Gateway key required).
func NewLambdaToolClient(url, apiKey string) *LambdaToolClient {
	return &LambdaToolClient{url: url, apiKey: apiKey}
}

// Call POSTs args as JSON to the tool endpoint and returns the parsed response body.
func (c *LambdaToolClient) Call(ctx context.Context, args map[string]any) (map[string]any, error) {
	payload, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("lambda tool client: marshal args: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("lambda tool client: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("lambda tool client: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("lambda tool client: read response: %w", err)
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("lambda tool client: parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lambda tool client: %s returned %d: %v", c.url, resp.StatusCode, result["error"])
	}

	return result, nil
}
