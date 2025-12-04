package usecase

import "context"

func (u *UseCase) DeleteLink(ctx context.Context, chatID int64, alias string) error {
	err := u.db.DeleteUserLink(ctx, chatID, alias)
	if err != nil {
		u.log.Error("DeleteLink: " + "Can't delete link" + err.Error())
		return err
	}

	return nil
}
