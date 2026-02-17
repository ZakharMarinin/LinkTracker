package usecase

import (
	"context"
	"errors"
	"scrapper/internal/domain"
)

func (u *UseCase) GetFilteredLinks(ctx context.Context, chatID int64, tags string) ([]*domain.Link, error) {
	const op = "UseCase::GetFilteredLinks"

	if chatID == 0 {
		return nil, errors.New("invalid chat id")
	}

	if tags == "" {
		return nil, errors.New("invalid tags")
	}

	links, err := u.DB.GetUserLinksByTag(ctx, chatID, tags)
	if err != nil {
		u.Log.Error("GetFilteredLinks: err with taking links: ", "op", op, "err", err)
		return nil, err
	}

	return links, nil
}
