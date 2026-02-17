package usecase_test

import (
	"context"
	"log/slog"
	"os"
	"scrapper/internal/config"
	"scrapper/internal/domain"
	"scrapper/internal/usecase"
	mocks "scrapper/internal/usecase/mocks"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCase_AddLink(t *testing.T) {
	type fields struct {
		db  *mocks.MockPostgres
		log *slog.Logger
		cfg *config.Config
	}

	type args struct {
		ctx    context.Context
		chatID int64
		url    string
		desc   string
		tags   string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				chatID: 1,
				url:    "https://google.com",
				desc:   "",
				tags:   "",
			},
		},
		{
			name: "no url test",
			args: args{
				ctx:    context.Background(),
				chatID: 2,
				url:    "",
				desc:   "",
				tags:   "",
			},
			wantErr: true,
		},
		{
			name: "no id test",
			args: args{
				ctx:  context.Background(),
				url:  "https://google.com",
				desc: "",
				tags: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostgres := mocks.NewMockPostgres(t)
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))
			cfg := &config.Config{}

			urlParts := strings.Split(tt.args.url, "/")
			alias := urlParts[len(urlParts)-1]

			mockPostgres.
				On("IsLinkExists", tt.args.ctx, tt.args.url).
				Maybe().
				Return(false, nil)

			mockPostgres.
				On("AddLink", tt.args.ctx, &domain.Link{ChatID: tt.args.chatID, URL: tt.args.url, Alias: alias}).
				Maybe().
				Return(nil)

			mockPostgres.
				On("IsUserLinkExists", tt.args.ctx, alias, tt.args.chatID).
				Maybe().
				Return(false, nil)

			mockPostgres.
				On("GetLinkByURL", tt.args.ctx, tt.args.url).
				Maybe().
				Return(&domain.Link{ChatID: tt.args.chatID, URL: tt.args.url, Alias: alias}, nil)

			mockPostgres.
				On("AddUserLink", tt.args.ctx, tt.args.chatID, &domain.Link{ChatID: tt.args.chatID, URL: tt.args.url, Alias: alias}).
				Maybe().
				Return(nil)

			u := &usecase.UseCase{
				DB:  mockPostgres,
				Log: log,
				Cfg: cfg,
			}

			err := u.AddLink(tt.args.ctx, tt.args.chatID, tt.args.url, tt.args.desc, tt.args.tags)
			if err != nil && !tt.wantErr {
				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
}

func TestUseCase_ExistedLink(t *testing.T) {
	type fields struct {
		db  usecase.Postgres
		log *slog.Logger
		cfg *config.Config
	}

	type args struct {
		ctx    context.Context
		chatID int64
		url    string
		desc   string
		tags   string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "already exists",
			args: args{
				ctx:    context.Background(),
				chatID: 1,
				url:    "https://google.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostgres := mocks.NewMockPostgres(t)
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))
			cfg := &config.Config{}

			urlParts := strings.Split(tt.args.url, "/")
			alias := urlParts[len(urlParts)-1]

			mockPostgres.
				On("IsLinkExists", tt.args.ctx, tt.args.url).
				Maybe().
				Return(true, nil)

			mockPostgres.
				On("AddLink", tt.args.ctx, &domain.Link{ChatID: tt.args.chatID, URL: tt.args.url, Alias: alias}).
				Maybe().
				Return(nil)

			mockPostgres.
				On("IsUserLinkExists", tt.args.ctx, alias, tt.args.chatID).
				Maybe().
				Return(true, nil)

			mockPostgres.
				On("GetLinkByURL", tt.args.ctx, tt.args.url).
				Maybe().
				Return(&domain.Link{ChatID: tt.args.chatID, URL: tt.args.url, Alias: alias}, nil)

			mockPostgres.
				On("AddUserLink", tt.args.ctx, tt.args.chatID, &domain.Link{ChatID: tt.args.chatID, URL: tt.args.url, Alias: alias}).
				Maybe().
				Return(nil)

			u := &usecase.UseCase{
				DB:  mockPostgres,
				Log: log,
				Cfg: cfg,
			}

			err := u.AddLink(tt.args.ctx, tt.args.chatID, tt.args.url, tt.args.desc, tt.args.tags)
			if err != nil && !tt.wantErr {
				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
}
