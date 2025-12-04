package cronModule

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"scrapper/internal/domain"
	"strings"
	"sync"

	"github.com/go-co-op/gocron"
)

type Storage interface {
	GetLinksToCheck(ctx context.Context, limit, offset uint64) ([]domain.Link, error)
	UpdateLink(ctx context.Context, link *domain.Link) (*domain.Link, error)
	DeleteLink(ctx context.Context, link *domain.Link) error
}

type GithubClient interface {
	GetUpdates(ctx context.Context, link *domain.Link) (update *domain.GitHubContent, err error)
}

type TGClient interface {
	SendUpdate(ctx context.Context, link *domain.Response) (*http.Response, error)
}

type Cron struct {
	log          *slog.Logger
	Cron         *gocron.Scheduler
	Storage      Storage
	GitHubClient GithubClient
	TGClient     TGClient
	Limit        uint64
}

func New(log *slog.Logger, cron *gocron.Scheduler, storage Storage, github GithubClient, tgClient TGClient, limit uint64) *Cron {
	return &Cron{
		log:          log,
		Cron:         cron,
		Storage:      storage,
		GitHubClient: github,
		TGClient:     tgClient,
		Limit:        limit,
	}
}

func (c *Cron) StartCron() {
	var offset uint64
	ctx := context.Background()

	c.log.Info("starting cron job")

	for {
		links, err := c.Storage.GetLinksToCheck(ctx, c.Limit, offset)
		if err != nil {
			c.log.Error("failed to fetch links", "error", err)
			break
		}

		if len(links) == 0 {
			break
		}

		c.StartWorkers(ctx, links)

		offset += c.Limit
		c.log.Info("cron job finished", "offset", offset)
	}
}

func (c *Cron) StartWorkers(ctx context.Context, links []domain.Link) {
	wg := &sync.WaitGroup{}
	workerCount := 10

	jobs := make(chan domain.Link)

	for i := 0; i < workerCount; i++ {
		go func() {
			for link := range jobs {
				err := c.ProcessLink(ctx, &link)
				if err != nil {
					c.log.Error("failed to process link", "error", err)
				}
				wg.Done()
			}
		}()
	}

	for _, link := range links {
		wg.Add(1)
		jobs <- link
	}

	wg.Wait()
	close(jobs)

	c.log.Info("cron job finished")
}

func (c *Cron) ProcessLink(ctx context.Context, link *domain.Link) error {
	var message string

	if strings.Contains(link.URL, "https://github.com") {
		update, err := c.GitHubClient.GetUpdates(ctx, link)
		if err != nil {
			return err
		}
		if update.PullRequests != nil || update.Issues != nil {
			message = CreateUpdateText(update)
		} else {
			return nil
		}
	} else {
		c.log.Info("invalid link", "url", link.URL)

		err := c.Storage.DeleteLink(ctx, link)
		if err != nil {
			c.log.Error("failed to delete link", "error", err)
			return err
		}

		return nil
	}

	_, err := c.Storage.UpdateLink(ctx, link)
	if err != nil {
		c.log.Error("failed to update link", "error", err)
		return err
	}

	err = c.SendUpdate(ctx, link, message)
	if err != nil {
		c.log.Error("failed to send update", "error", err)
		return err
	}

	c.log.Info("processed link", "url", link.URL)
	return nil
}

func (c *Cron) SendUpdate(ctx context.Context, link *domain.Link, message string) error {
	linkUpdate := &domain.Link{
		URL:  link.URL,
		Desc: message,
	}

	req := &domain.Response{
		ChatIDs: []int64{link.ChatID},
		Link:    linkUpdate,
	}

	_, err := c.TGClient.SendUpdate(ctx, req)
	if err != nil {
		c.log.Error("failed to send update", "error", err)
		return err
	}

	c.log.Info("sent update", "url", link.URL)

	return nil
}

func CreateUpdateText(repoData *domain.GitHubContent) string {
	message := ""

	for _, i := range repoData.Issues {
		formatTine := i.UpdatedAt.Format("2006-01-02 15:04")
		formatUser := strings.Replace(fmt.Sprint(i.User), "{", "", -1)
		formatUser = strings.Replace(formatUser, "}", "", -1)
		message += fmt.Sprintf("Тип: Issue\nЗаголовок: %s\nОписание: %s\nКем: %s\nКогда: %s\n\n", i.Title, i.Body, formatUser, formatTine)
	}

	for _, i := range repoData.PullRequests {
		formatTine := i.UpdatedAt.Format("2006-01-02 15:04")
		formatUser := strings.Replace(fmt.Sprint(i.User), "{", "", -1)
		formatUser = strings.Replace(formatUser, "}", "", -1)
		message += fmt.Sprintf("Тип: PullRequest\nЗаголовок: %s\nОписание: %s\nКем: %s\nКогда: %s\n\n", i.Title, i.Body, formatUser, formatTine)
	}

	return message
}
