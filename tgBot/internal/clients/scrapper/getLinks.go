package scrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"linktracker/internal/domain"
	"net/http"
)

func (s *Client) GetLinks(ctx context.Context, chatID int64) ([]*domain.Link, error) {
	url := fmt.Sprintf("%s/links/%d", s.addr, chatID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		s.log.Error("GetLinks: Error creating request to tg-chat", "error", err)
		return nil, err
	}

	resp, err := s.sendRequest(req)
	if err != nil {
		s.log.Error("GetLinks: Error sending request to tg-chat", "error", err)
		return nil, err
	}

	var links []*domain.Link

	err = json.Unmarshal(resp, &links)
	if err != nil {
		s.log.Error("GetLinks: Error unmarshalling response", "error", err)
		return nil, err
	}

	return links, nil
}
