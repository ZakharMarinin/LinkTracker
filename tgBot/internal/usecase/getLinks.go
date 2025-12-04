package usecase

import (
	"context"
	"linktracker/internal/domain"
)

func (u *UseCase) GetLinks(ctx context.Context, id int64) ([]*domain.Link, error) {
	u.log.Info("GetLinks: ready to get links")

	resp, err := u.ScrapperClient.GetLinks(ctx, id)
	if err != nil {
		u.log.Error("GetLinks: Error getting links", "error", err)
		return nil, err
	}

	return resp, nil
}
