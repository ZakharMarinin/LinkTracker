package usecase

import (
	"context"
)

func (u *UseCase) DeleteLink(ctx context.Context, id int64, alias string) error {
	u.log.Info("Starting DeleteLink")

	err := u.ScrapperClient.DeleteLink(ctx, id, alias)
	if err != nil {
		u.log.Error(err.Error())
		return err
	}

	return nil
}
