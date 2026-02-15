package usecase

import (
	"context"
	"errors"
)

func (u *UseCase) DeleteLink(ctx context.Context, chatID int64, alias string) error {
	if chatID == 0 {
		return errors.New("invalid chat id")
	}

	if alias == "" {
		return errors.New("invalid alias")
	}

	err := u.db.DeleteUserLink(ctx, chatID, alias)
	if err != nil {
		u.log.Error("DeleteLink: " + "Can't delete link" + err.Error())
		return err
	}

	return nil
}
