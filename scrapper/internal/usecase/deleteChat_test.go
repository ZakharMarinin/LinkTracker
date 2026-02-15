package usecase

import (
	"context"
	"log/slog"
	"os"
	"scrapper/internal/config"
	mocks "scrapper/internal/usecase/mocks"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCase_DeleteChat(t *testing.T) {
	type fields struct {
		db  Postgres
		log *slog.Logger
		cfg *config.Config
	}
	type args struct {
		ctx    context.Context
		chatID int64
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
				chatID: int64(1),
			},
			wantErr: false,
		},
		{
			name: "invalid chatID",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPostgres := mocks.NewMockPostgres(t)
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))
			cfg := &config.Config{}

			mockPostgres.
				On("DeleteChat", tt.args.ctx, tt.args.chatID).
				Maybe().
				Return(nil)

			u := &UseCase{
				db:  mockPostgres,
				log: log,
				cfg: cfg,
			}
			err := u.DeleteChat(tt.args.ctx, tt.args.chatID)
			if err != nil && !tt.wantErr {
				t.Errorf("DeleteChat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
}
