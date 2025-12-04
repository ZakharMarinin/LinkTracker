package scrapper

import (
	"fmt"
	"io"
	"linktracker/internal/config"
	"linktracker/internal/domain"
	"log/slog"
	"net/http"
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

func NewScrapperClient(log *slog.Logger, cfg *config.Config) *Client {
	httpClient := &http.Client{
		Timeout: cfg.BotClients.Scrapper.Timeout,
	}
	return &Client{addr: cfg.BotClients.Scrapper.Addr, client: httpClient, log: log}
}

// Sending Request
func (s *Client) sendRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		s.log.Error("Error sending request", "url", req.URL.String(), "method", req.Method, "error", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		s.log.Error("Error executing request", "status", res.Status)
		return nil, fmt.Errorf(res.Status)
	}

	return io.ReadAll(res.Body)
}
