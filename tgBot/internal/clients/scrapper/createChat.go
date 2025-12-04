package scrapper

import (
	"context"
	"fmt"
	"net/http"
)

func (s *Client) CreateChat(ctx context.Context, chatID int64) error {
	url := fmt.Sprintf("%s/tg-chat/%d", s.addr, chatID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		s.log.Error("CreateChat: Error creating request to tg-chat: " + err.Error())
		return err
	}

	resp, err := s.sendRequest(req)
	if err != nil {
		s.log.Error("CreateChat: Error sending request to tg-chat: " + err.Error())
		return err
	}

	s.log.Info("CreateChat: Successfully sent request to tg-chat: " + fmt.Sprint(resp))

	return nil
}
