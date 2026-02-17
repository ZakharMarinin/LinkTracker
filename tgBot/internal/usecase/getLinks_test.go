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
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCase_GetLinks(t *testing.T) {
	type fields struct {
		log            *slog.Logger
		ScrapperClient usecase.ScrapperClient
		Storage        usecase.Storage
	}

	type args struct {
		ctx context.Context
		id  int64
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
				ctx: context.Background(),
				id:  int64(1),
			},
			want:    []*domain.Link{},
			wantErr: false,
		},
		{
			name: "error getting temp user",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockScrapperClient(t)
			mockStorage := mocks.NewMockStorage(t)
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))

			mockStorage.
				On("GetTempUserLinks", tt.args.ctx, tt.args.id).
				Maybe().
				Return(func(chatID int64) (*storage.TempUserLinks, error) {
					if chatID != 0 {
						return &storage.TempUserLinks{UserID: tt.args.id, Links: []*domain.Link{}}, nil
					}

					return nil, errors.New("invalid id")
				}(tt.args.id))

			mockClient.
				On("GetLinks", tt.args.ctx, tt.args.id).
				Maybe().
				Return(func(id int64) ([]*domain.Link, error) {
					if id != 0 {
						return []*domain.Link{}, nil
					}

					return nil, errors.New("invalid id")
				}(tt.args.id))

			mockStorage.
				On("SaveTempUserLinks", tt.args.ctx, &storage.TempUserLinks{
					UserID: tt.args.id,
					Links:  []*domain.Link{},
				}).
				Maybe().
				Return(nil)

			u := &usecase.UseCase{
				Log:            log,
				ScrapperClient: mockClient,
				Storage:        mockStorage,
			}

			got, err := u.GetLinks(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLinks() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				require.Error(t, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinks() got = %v, want %v", got, tt.want)
			}
		})
	}
}
