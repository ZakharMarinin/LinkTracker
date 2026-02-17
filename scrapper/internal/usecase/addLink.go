package usecase

import (
	"context"
	"errors"
	"fmt"
	"scrapper/internal/domain"
	"strings"
)

func (u *UseCase) AddLink(ctx context.Context, chatID int64, url, desc, tags string) error {
	if chatID == 0 {
		return errors.New("invalid chat id")
	}

	if url == "" {
		return errors.New("invalid url")
	}

	exists, err := u.DB.IsLinkExists(ctx, url)
	if err != nil {
		u.Log.Error("Error while checking if link exists", "error", err)
		return err
	}

	urlParts := strings.Split(url, "/")
	alias := urlParts[len(urlParts)-1]

	link := &domain.Link{
		ChatID: chatID,
		URL:    url,
		Desc:   desc,
		Tags:   tags,
		Alias:  alias,
	}

	if !exists {
		err := u.DB.AddLink(ctx, link)
		if err != nil {
			u.Log.Error("AddLink", "error", err)
			return err
		}
	}

	isIt, err := u.DB.IsUserLinkExists(ctx, link.Alias, link.ChatID)
	if err != nil {
		u.Log.Error("AddLink", "error", err)
		return err
	}

	if !isIt {
		existsLink, err := u.DB.GetLinkByURL(ctx, url)
		if err != nil {
			u.Log.Error("GetLinkByAlias", "error", err)
			return err
		}

		link.ID = existsLink.ID

		err = u.DB.AddUserLink(ctx, chatID, link)
		if err != nil {
			u.Log.Error("AddUserLink", "error", err)
			return err
		}

		return nil
	} else {
		return fmt.Errorf("already link exists")
	}
}
