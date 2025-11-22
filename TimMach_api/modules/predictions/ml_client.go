package predictions

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// RestyMLClient là impl MLClient dùng resty (có timeout, base URL).
type RestyMLClient struct {
	client *resty.Client
}

func NewRestyMLClient(baseURL string, timeout time.Duration) *RestyMLClient {
	return &RestyMLClient{
		client: resty.New().SetBaseURL(baseURL).SetTimeout(timeout),
	}
}

func (c *RestyMLClient) Predict(ctx context.Context, payload MLRequest) (MLResponse, error) {
	if c == nil || c.client == nil {
		return MLResponse{}, fmt.Errorf("ml client is not configured")
	}
	var out MLResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&out).
		Post("/predict")
	if err != nil {
		return MLResponse{}, err
	}
	if resp.IsError() {
		return MLResponse{}, fmt.Errorf("ml service status %d", resp.StatusCode())
	}
	return out, nil
}
