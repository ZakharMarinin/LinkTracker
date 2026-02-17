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

	err := u.DB.DeleteUserLink(ctx, chatID, alias)
	if err != nil {
		u.Log.Error("DeleteLink: " + "Can't delete link" + err.Error())
		return err
	}

	return nil
}
