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

	go func() {
		u.log.Info("AddLink: Cache invalidating link", "link", link)

		links, err := u.Storage.GetTempUserLinks(ctx, id)
		if err != nil {
			u.log.Error("AddLink: cannot get cache", "err", err)
		}

		if links != nil {
			links.Links = append(links.Links, &link)

			err = u.Storage.SaveTempUserLinks(ctx, links)
			if err != nil {
				u.log.Error("AddLink: cannot save cache", "err", err)
			}
		}
	}()

	return nil
}
