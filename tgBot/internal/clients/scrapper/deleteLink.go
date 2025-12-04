package scrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"linktracker/internal/domain"
	"net/http"
)

func (s *Client) DeleteLink(ctx context.Context, chatID int64, alias string) error {
	var body domain.DeleteLinkRequest

	body.ChatID = chatID
	body.Alias = alias

	bodyReq, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", s.addr+"/links", bytes.NewReader(bodyReq))
	if err != nil {
		s.log.Error("DeleteLink: Error creating request to tg-chat: " + err.Error())
		return err
	}

	_, err = s.sendRequest(req)
	if err != nil {
		s.log.Error("DeleteLink: Error sending request to tg-chat: " + err.Error())
		return err
	}

	return nil
}
