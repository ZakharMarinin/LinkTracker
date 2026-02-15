package usecase

import (
	"context"
	"errors"
	"scrapper/internal/domain"
)

func (u *UseCase) GetLinks(ctx context.Context, chatID int64) ([]domain.Link, error) {
	if chatID == 0 {
		return nil, errors.New("invalid chat id")
	}

	userLinks, err := u.db.GetLinksByChatID(ctx, chatID)
	if err != nil {
		u.log.Error("UseCase-GetLinks: err with taking links: ", "error", err)
		return nil, err
	}

	return userLinks, nil
}
