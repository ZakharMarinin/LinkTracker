package usecase

import (
	"context"
	"linktracker/internal/domain"
	"linktracker/internal/storage"
)

func (u *UseCase) GetLinks(ctx context.Context, id int64) ([]*domain.Link, error) {
	u.log.Info("GetLinks: ready to get links")

	isCached := true

	links, err := u.Storage.GetTempUserLinks(ctx, id)
	if err != nil {
		isCached = false
		u.log.Error("GetLinks: error getting temp user links: ", err)
	}

	if links == nil {
		isCached = false
		u.log.Info("GetLinks: no temp user links")
	}

	if !isCached {
		resp, err := u.ScrapperClient.GetLinks(ctx, id)
		if err != nil {
			u.log.Error("GetLinks: Error getting links", "error", err)
			return nil, err
		}

		tempUserLinks := &storage.TempUserLinks{
			UserID: id,
			Links:  resp,
		}

		err = u.Storage.SaveTempUserLinks(ctx, tempUserLinks)
		if err != nil {
			u.log.Error("GetLinks: Error saving temp user links: ", err)
		}

		return resp, nil
	}

	u.log.Info("GetLinks: Cache validating link successfully")
	
	return links.Links, nil
}
