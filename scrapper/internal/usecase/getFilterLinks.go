package usecase

import (
	"context"
	"scrapper/internal/domain"
)

func (u *UseCase) GetFilteredLinks(ctx context.Context, chatID int64, tags string) ([]*domain.Link, error) {
	const op = "UseCase::GetFilteredLinks"

	links, err := u.db.GetUserLinksByTag(ctx, chatID, tags)
	if err != nil {
		u.log.Error("GetFilteredLinks: err with taking links: ", "op", op, "err", err)
		return nil, err
	}

	return links, nil
}
