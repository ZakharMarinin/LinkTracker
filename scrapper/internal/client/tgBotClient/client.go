package tgBotClient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"scrapper/internal/config"
	"scrapper/internal/domain"
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

func (c *Client) doRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		c.log.Error("cannot create request: ", err)
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.log.Error("cannot send request: ", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) SendUpdate(ctx context.Context, link *domain.Response) (*http.Response, error) {
	url := c.addr + "/updates"

	body, err := json.Marshal(link)
	if err != nil {
		c.log.Error("cannot marshal update request: ", err)
		return nil, err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		c.log.Error("cannot send request: ", err)
		return nil, err
	}

	return resp, nil
}
