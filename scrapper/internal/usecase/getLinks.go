package usecase

import (
	"context"
	"scrapper/internal/domain"
)

func (u *UseCase) GetLinks(ctx context.Context, chatID int64) ([]domain.Link, error) {
	userLinks, err := u.db.GetLinksByChatID(ctx, chatID)
	if err != nil {
		u.log.Error("UseCase-GetLinks: err with taking links: ", err)
		return nil, err
	}

	return userLinks, nil
}
