package scrapper

import (
	"bytes"
	"fmt"
	"io"
	"linktracker/internal/config"
	"linktracker/internal/domain"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type Response struct {
	ChatID int64          `json:"chat_id"`
	Link   []*domain.Link `json:"link"`
	Desc   string         `json:"desc"`
	Msg    string         `json:"msg"`
}

type Client struct {
	addr   string
	client *http.Client
	log    *slog.Logger
}

const (
	maxRetries = 5
	baseDelay  = time.Second * 1
	maxDelay   = time.Second * 5
)

func NewScrapperClient(log *slog.Logger, cfg *config.Config) *Client {
	httpClient := &http.Client{
		Timeout: cfg.BotClients.Scrapper.Timeout,
	}

	return &Client{addr: cfg.BotClients.Scrapper.Addr, client: httpClient, log: log}
}

func (s *Client) sendRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := s.DoWithRetry(req)
	if err != nil {
		s.log.Error("Error sending request", "url", req.URL.String(), "method", req.Method, "error", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		s.log.Error("Error executing request", "status", res.Status)
		return nil, fmt.Errorf("%s", res.Status)
	}

	return io.ReadAll(res.Body)
}

func (s *Client) DoWithRetry(req *http.Request) (*http.Response, error) {
	var bodyData []byte

	if req.Body != nil {
		var err error

		bodyData, err = io.ReadAll(req.Body)
		if err != nil {
			s.log.Error("cannot read request body", "error", err)

			return nil, err
		}

		req.Body.Close()
	}

	for i := 1; i <= maxRetries; i++ {
		if bodyData != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyData))
		}

		resp, err := s.client.Do(req)

		shouldRetry := false

		if err != nil {
			shouldRetry = true

			s.log.Error("cannot send request", "error", err, "attempt", i, "maxRetries", maxRetries)
		} else if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			shouldRetry = true

			s.log.Error("cannot send request", "error", err, "attempt", i, "maxRetries", maxRetries)
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
