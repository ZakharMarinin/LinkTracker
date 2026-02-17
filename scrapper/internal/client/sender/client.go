package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"scrapper/internal/config"
	"scrapper/internal/domain"
	"time"
)

type Client struct {
	addr   string
	client *http.Client
	log    *slog.Logger
}

func NewTGClient(log *slog.Logger, cfg *config.Config) *Client {
	httpClient := &http.Client{
		Timeout: cfg.TgBot.Timeout,
	}

	return &Client{addr: cfg.TgBot.Addr, client: httpClient, log: log}
}

func (c *Client) doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		c.log.Error("cannot create request", "error", err)

		return nil, err
	}

	resp, err := c.DoWithRetry(req)
	if err != nil {
		c.log.Error("cannot send request", "error", err)

		return nil, err
	}

	return resp, nil
}

func (c *Client) DoWithRetry(req *http.Request) (*http.Response, error) {
	const (
		maxRetries = 5
		baseDelay  = time.Second * 1
		maxDelay   = time.Second * 5
	)

	var bodyData []byte

	if req.Body != nil {
		var err error

		bodyData, err = io.ReadAll(req.Body)
		if err != nil {
			c.log.Error("cannot read request body", "error", err)

			return nil, err
		}

		req.Body.Close()
	}

	for i := 1; i <= maxRetries; i++ {
		if bodyData != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyData))
		}

		resp, err := c.client.Do(req)

		shouldRetry := false

		if err != nil {
			shouldRetry = true

			c.log.Error("cannot send request", "error", err, "attempt", i, "maxRetries", maxRetries)
		} else if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			shouldRetry = true

			c.log.Error("cannot send request", "error", err, "attempt", i, "maxRetries", maxRetries)
			resp.Body.Close()
		}

		if !shouldRetry {
			return resp, nil
		}

		delay := float64(baseDelay) * math.Pow(2, float64(i))
		jitter := (rand.Float64() * 0.2) + 0.9
		sleepDuration := time.Duration(delay * jitter)

		if sleepDuration > maxDelay {
			sleepDuration = maxDelay
		}

		select {
		case <-time.After(sleepDuration):
			continue
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

func (c *Client) SendUpdate(ctx context.Context, link *domain.Response) error {
	url := c.addr + "/updates"

	body, err := json.Marshal(link)
	if err != nil {
		c.log.Error("cannot marshal update request", "error", err)

		return err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		c.log.Error("cannot send request", "error", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
