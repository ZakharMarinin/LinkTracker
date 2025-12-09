package scrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"linktracker/internal/domain"
	"net/http"
)

func (s *Client) GetFilteredLinks(ctx context.Context, chatID int64, tag string) ([]*domain.Link, error) {
	const op = "Client::GetFilteredLinks"

	url := fmt.Sprintf("%s/links/%d/%s", s.addr, chatID, tag)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		s.log.Error("Error creating request to tg-chat: ", "op", op, "error", err.Error())
		return nil, err
	}

	resp, err := s.sendRequest(req)
	if err != nil {
		s.log.Error("Error sending request to tg-chat: ", "op", op, "error", err.Error())
		return nil, err
	}

	var links []*domain.Link
	err = json.Unmarshal(resp, &links)
	if err != nil {
		s.log.Error("Error unmarshalling response: ", "op", op, "error", err.Error())
		return nil, err
	}

	return links, nil
}
