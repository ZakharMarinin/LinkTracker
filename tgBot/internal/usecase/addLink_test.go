package usecase_test

import (
	"context"
	"errors"
	"linktracker/internal/domain"
	"linktracker/internal/storage"
	"linktracker/internal/usecase"
	mocks "linktracker/internal/usecase/mocks"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCase_AddLink(t *testing.T) {
	type fields struct {
		log            *slog.Logger
		ScrapperClient usecase.ScrapperClient
		Storage        usecase.Storage
	}

	type args struct {
		ctx  context.Context
		id   int64
		link *domain.Link
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  int64(1),
				link: &domain.Link{
					URL:    "https://google.com",
					Domain: "google.com",
					Desc:   "asdfasdf",
					Tags:   "asdf, asdf, asd",
					ChatID: int64(1),
				},
			},
			wantErr: false,
		},
		{
			name: "Client error",
			args: args{
				ctx:  context.Background(),
				id:   int64(1),
				link: &domain.Link{},
			},
			wantErr: true,
		},
		{
			name: "Cannot get cache error",
			args: args{
				ctx: context.Background(),
				link: &domain.Link{
					URL:    "https://google.com",
					Domain: "google.com",
					Desc:   "asdfasdf",
					Tags:   "asdf, asdf, asd",
				},
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
				On("AddLink", tt.args.ctx, tt.args.id, tt.args.link).
				Maybe().
				Return(func(link domain.Link) error {
					if link.URL != "" {
						return nil
					}

					return errors.New("invalid link")
				}(*tt.args.link))

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
					Links: []*domain.Link{
						tt.args.link,
					},
				}).
				Maybe().
				Return(nil)

			u := &usecase.UseCase{
				Log:            log,
				ScrapperClient: mockClient,
				Storage:        mockStorage,
			}

			err := u.AddLink(tt.args.ctx, tt.args.id, tt.args.link)
			if err != nil && !tt.wantErr {
				t.Errorf("AddLink() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
}
