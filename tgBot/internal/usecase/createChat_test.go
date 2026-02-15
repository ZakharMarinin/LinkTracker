package usecase

import (
	"context"
	"errors"
	mocks "linktracker/internal/usecase/mocks"
	"log/slog"
	"os"
	"testing"
)

func TestUseCase_CreateChat(t *testing.T) {
	type fields struct {
		log            *slog.Logger
		ScrapperClient ScrapperClient
		Storage        Storage
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
			name: "Success",
			args: args{
				ctx:    context.Background(),
				chatID: int64(1),
			},
			wantErr: false,
		},
		{
			name: "Invalid id error",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockScrapperClient(t)
			mockStorage := mocks.NewMockStorage(t)
			log := slog.New(slog.NewTextHandler(os.Stdout, nil))

			mockClient.
				On("CreateChat", tt.args.ctx, tt.args.chatID).
				Once().
				Return(func(chatID int64) error {
					if chatID != 0 {
						return nil
					}
					return errors.New("invalid id")
				}(tt.args.chatID))

			u := &UseCase{
				log:            log,
				ScrapperClient: mockClient,
				Storage:        mockStorage,
			}
			if err := u.CreateChat(tt.args.ctx, tt.args.chatID); (err != nil) != tt.wantErr {
				t.Errorf("CreateChat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
