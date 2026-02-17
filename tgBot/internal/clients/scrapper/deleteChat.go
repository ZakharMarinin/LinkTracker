package scrapper

import (
	"context"
	"fmt"
	"net/http"
)

func (s *Client) DeleteChat(ctx context.Context, chatID int64) error {
	url := fmt.Sprintf("%s/links/%d", s.addr, chatID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, http.NoBody)
	if err != nil {
		s.log.Error("DeleteChat: Error creating request to tg-chat", "error", err)
		return err
	}

	_, err = s.sendRequest(req)
	if err != nil {
		s.log.Error("DeleteChat: Error sending request to tg-chat", "error", err)
		return err
	}

	return nil
}
