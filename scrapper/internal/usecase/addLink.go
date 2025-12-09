package usecase

import (
	"context"
	"fmt"
	"scrapper/internal/domain"
	"strings"
)

func (u *UseCase) AddLink(ctx context.Context, chatID int64, url, desc, tags string) error {
	exists, err := u.db.IsLinkExists(ctx, url)
	if err != nil {
		u.log.Error("Error while checking if link exists", "error", err)
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
		err := u.db.AddLink(ctx, link)
		if err != nil {
			u.log.Error("AddLink", "error", err)
			return err
		}
	}

	isIt, err := u.db.IsUserLinkExists(ctx, link.Alias, link.ChatID)
	if err != nil {
		u.log.Error("AddLink", "error", err)
		return err
	}

	if !isIt {
		existsLink, err := u.db.GetLinkByURL(ctx, url)
		if err != nil {
			u.log.Error("GetLinkByAlias", "error", err)
			return err
		}

		link.ID = existsLink.ID

		err = u.db.AddUserLink(ctx, chatID, link)
		if err != nil {
			u.log.Error("AddUserLink", "error", err)
			return err
		}

		return nil
	} else {
		return fmt.Errorf("already link exists")
	}
}
