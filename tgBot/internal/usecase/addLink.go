package usecase

import (
	"context"
	"linktracker/internal/domain"
)

func (u *UseCase) AddLink(ctx context.Context, id int64, link domain.Link) error {
	u.log.Info("ready to add link")

	err := u.ScrapperClient.AddLink(ctx, id, link)
	if err != nil {
		u.log.Error("AddLink: cannot s%w", err)
		return err
	}

	return nil
}
