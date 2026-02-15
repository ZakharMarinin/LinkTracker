package usecase

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"scrapper/internal/config"
	"scrapper/internal/domain"
	mocks "scrapper/internal/usecase/mocks"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCase_GetFilteredLinks(t *testing.T) {
	type fields struct {
		db  Postgres
		log *slog.Logger
		cfg *config.Config
	}
	type args struct {
		ctx    context.Context
		chatID int64
		tags   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domain.Link
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				chatID: int64(1),
				tags:   "test",
			},
			want:    []*domain.Link{},
			wantErr: false,
		},
		{
			name: "invalid tags",
			args: args{
				ctx:    context.Background(),
				chatID: int64(1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid chatID",
			args: args{
				ctx:  context.Background(),
				tags: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostgres := mocks.NewMockPostgres(t)
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))
			cfg := &config.Config{}

			mockPostgres.
				On("GetUserLinksByTag", tt.args.ctx, tt.args.chatID, tt.args.tags).
				Maybe().
				Return([]*domain.Link{}, nil)

			u := &UseCase{
				db:  mockPostgres,
				log: log,
				cfg: cfg,
			}
			got, err := u.GetFilteredLinks(tt.args.ctx, tt.args.chatID, tt.args.tags)
			if err != nil && !tt.wantErr {
				t.Errorf("GetFilteredLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.Error(t, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFilteredLinks() got = %v, want %v", got, tt.want)
			}
		})
	}
}
