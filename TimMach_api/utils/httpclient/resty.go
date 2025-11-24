package httpclient

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// NewRestyClient creates a resty client with base URL and timeout pre-set.
func NewRestyClient(baseURL string, timeout time.Duration) *resty.Client {
	return resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout)
}

// PostJSON sends a JSON POST request and unmarshals the response into out.
// Returns an error if the request fails or the service returns a non-2xx status.
func PostJSON(ctx context.Context, client *resty.Client, path string, payload any, out any) (*resty.Response, error) {
	if client == nil {
		return nil, fmt.Errorf("http client is not configured")
	}
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(out).
		Post(path)
	if err != nil {
		return resp, err
	}
	if resp.IsError() {
		return resp, fmt.Errorf("service status %d", resp.StatusCode())
	}
	return resp, nil
}
