package handlers

import (
	"context"
	"linktracker/internal/domain"
	"log/slog"
	"strings"

	"gopkg.in/telebot.v4"
)

type UseCase interface {
	AddLink(ctx context.Context, id int64, link domain.Link) error
	DeleteLink(ctx context.Context, id int64, alias string) error
	GetLinks(ctx context.Context, id int64) ([]*domain.Link, error)
	GetFilteredLinks(ctx context.Context, id int64, tag string) ([]*domain.Link, error)
	CreateChat(ctx context.Context, chatID int64) error
	ChangeUserState(ctx context.Context, userInfo *domain.UserStateInfo, state string) error
	GetUserState(ctx context.Context, userID int64) (*domain.UserStateInfo, error)
}

type BotHandler struct {
	Bot     *telebot.Bot
	useCase UseCase
	log     *slog.Logger
}

func NewBotHandler(b *telebot.Bot, useCase UseCase, log *slog.Logger) *BotHandler {
	return &BotHandler{b, useCase, log}
}

func (b *BotHandler) Cancel(ctx context.Context, userInfo *domain.UserStateInfo, c telebot.Context) error {
	err := b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
	if err != nil {
		b.log.Error("Failed to change user state", "userInfo", userInfo, "err", err)
		return err
	}

	b.Bot.Send(c.Recipient(), "Операция была прервана.")

	return nil
}

func (b *BotHandler) LinkValidation(url string) (*domain.Link, bool) {
	_, git, isIt := strings.Cut(url, "https://")
	if isIt {
		linkParts := strings.Split(git, "/")

		if len(linkParts) != 3 {
			return nil, false
		}

		if len(linkParts[0]) <= 0 || len(linkParts[1]) <= 0 || len(linkParts[2]) <= 0 {
			return nil, false
		}

		if linkParts[0] == "github.com" {
			link := &domain.Link{
				Domain:     linkParts[0],
				Author:     linkParts[1],
				Repository: linkParts[2],
			}
			return link, true
		}
	}
	return nil, false
}
