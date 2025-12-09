package usecase

import (
	"context"
	"strings"
)

func (u *UseCase) DeleteLink(ctx context.Context, id int64, alias string) error {
	u.log.Info("Starting DeleteLink")

	err := u.ScrapperClient.DeleteLink(ctx, id, alias)
	if err != nil {
		u.log.Error(err.Error())
		return err
	}

	go func() {
		u.log.Info("DeleteLink: Cache invalidating link")

		links, err := u.Storage.GetTempUserLinks(ctx, id)
		if err != nil {
			u.log.Error("AddLink: cannot s%w", err)
		}

		if links != nil {
			for i := range links.Links {
				linkParts := strings.Split(links.Links[i].URL, "/")

				if linkParts[len(linkParts)-1] == alias {
					links.Links = append(links.Links[:i], links.Links[i+1:]...)
				}
			}

			err = u.Storage.SaveTempUserLinks(ctx, links)
			if err != nil {
				u.log.Error("AddLink: cannot s%w", err)
			}
		}
	}()

	return nil
}
