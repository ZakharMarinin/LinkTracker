package scrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"linktracker/internal/domain"
	"net/http"
)

func (s *Client) AddLink(ctx context.Context, chatID int64, link domain.Link) error {
	body, err := json.Marshal(&link)
	if err != nil {
		s.log.Error("AddLink: Error marshalling link: " + err.Error())
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.addr+"/links", bytes.NewReader(body))
	if err != nil {
		s.log.Error("AddLink: Error creating request to tg-chat: " + err.Error())
		return err
	}

	req.Header.Set("Chat-ID", fmt.Sprint(chatID))

	_, err = s.sendRequest(req)
	if err != nil {
		s.log.Error("AddLink: Error sending request to tg-chat: " + err.Error())
		return err
	}

	return nil
}
