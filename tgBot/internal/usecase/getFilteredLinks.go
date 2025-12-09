package usecase

import (
	"context"
	"linktracker/internal/domain"
	"strings"
)

func (u *UseCase) GetFilteredLinks(ctx context.Context, id int64, tag string) ([]*domain.Link, error) {
	const op = "UseCase::GetFilteredLinks"

	u.log.Info("Ready to get links", "op", op)

	tag, _ = strings.CutPrefix(tag, "#")

	resp, err := u.ScrapperClient.GetFilteredLinks(ctx, id, tag)
	if err != nil {
		u.log.Error("Error getting links", "op", op, "error", err)
		return nil, err
	}

	return resp, nil
}
