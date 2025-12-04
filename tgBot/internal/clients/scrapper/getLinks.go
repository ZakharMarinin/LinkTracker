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

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		s.log.Error("GetLinks: Error creating request to tg-chat: " + err.Error())
		return nil, err
	}

	resp, err := s.sendRequest(req)
	if err != nil {
		s.log.Error("GetLinks: Error sending request to tg-chat: " + err.Error())
		return nil, err
	}

	var links []*domain.Link
	err = json.Unmarshal(resp, &links)
	if err != nil {
		s.log.Error("GetLinks: Error unmarshalling response: " + err.Error())
		return nil, err
	}

	return links, nil
}
