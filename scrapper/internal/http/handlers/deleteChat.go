package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *HTTP) DeleteChat(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := chi.URLParam(r, "id")

		if chatID == "" {
			h.log.Error("Empty chat ID")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		intChatID, err := strconv.ParseInt(chatID, 10, 64)
		if err != nil {
			h.log.Error("could not parse int chatID")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.useCase.DeleteChat(ctx, intChatID)
		if err != nil {
			h.log.Error("could not delete chat")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.log.Info("successfully deleted chat")
		w.WriteHeader(http.StatusOK)
	}
}
