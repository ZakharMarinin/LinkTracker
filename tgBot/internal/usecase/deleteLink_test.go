package usecase_test

import (
	"context"
	"errors"
	"linktracker/internal/storage"
	"linktracker/internal/usecase"
	mocks "linktracker/internal/usecase/mocks"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCase_DeleteLink(t *testing.T) {
	type fields struct {
		log            *slog.Logger
		ScrapperClient usecase.ScrapperClient
		Storage        usecase.Storage
	}

	type args struct {
		ctx   context.Context
		id    int64
		alias string
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
				ctx:   context.Background(),
				id:    int64(1),
				alias: "test",
			},
			wantErr: false,
		},
		{
			name: "invalid alias",
			args: args{
				ctx: context.Background(),
				id:  int64(1),
			},
			wantErr: true,
		},
		{
			name: "Id error",
			args: args{
				ctx:   context.Background(),
				alias: "test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockScrapperClient(t)
			mockStorage := mocks.NewMockStorage(t)
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))

			mockClient.
				On("DeleteLink", tt.args.ctx, tt.args.id, tt.args.alias).
				Once().
				Return(func(chatID int64, alias string) error {
					if alias != "" || chatID != tt.args.id {
						return nil
					}

					return errors.New("invalid id or alias")
				}(tt.args.id, tt.args.alias))

			mockStorage.
				On("GetTempUserLinks", tt.args.ctx, tt.args.id).
				Maybe().
				Return(func(chatID int64) (*storage.TempUserLinks, error) {
					if chatID != 0 {
						return &storage.TempUserLinks{UserID: tt.args.id}, nil
					}

					return nil, errors.New("invalid id")
				}(tt.args.id))

			mockStorage.
				On("SaveTempUserLinks", tt.args.ctx, &storage.TempUserLinks{
					UserID: tt.args.id,
					Links:  nil,
				}).
				Maybe().
				Return(nil)

			u := &usecase.UseCase{
				Log:            log,
				ScrapperClient: mockClient,
				Storage:        mockStorage,
			}

			err := u.DeleteLink(tt.args.ctx, tt.args.id, tt.args.alias)
			if err != nil && !tt.wantErr {
				t.Errorf("DeleteLink() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
}
