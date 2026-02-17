package usecase_test

import (
	"context"
	"errors"
	"linktracker/internal/domain"
	"linktracker/internal/usecase"
	mocks "linktracker/internal/usecase/mocks"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUseCase_GetFilteredLinks(t *testing.T) {
	type fields struct {
		log            *slog.Logger
		ScrapperClient usecase.ScrapperClient
		Storage        usecase.Storage
	}

	type args struct {
		ctx context.Context
		id  int64
		tag string
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
				tag: "test",
			},
			want:    []*domain.Link{},
			wantErr: false,
		},
		{
			name: "id error",
			args: args{
				ctx: context.Background(),
				tag: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "tag error",
			args: args{
				ctx: context.Background(),
				id:  int64(1),
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

			mockClient.
				On("GetFilteredLinks", tt.args.ctx, tt.args.id, tt.args.tag).
				Once().
				Return(func(id int64, tag string) ([]*domain.Link, error) {
					if id != 0 && tag != "" {
						return []*domain.Link{}, nil
					}

					return nil, errors.New("not found")
				}(tt.args.id, tt.args.tag))

			u := &usecase.UseCase{
				Log:            log,
				ScrapperClient: mockClient,
				Storage:        mockStorage,
			}

			got, err := u.GetFilteredLinks(tt.args.ctx, tt.args.id, tt.args.tag)
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
