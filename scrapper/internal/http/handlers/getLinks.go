package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *HTTP) GetLinks(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := chi.URLParam(r, "id")

		if chatID == "" {
			h.log.Error("handler-GetLinks: Empty chat ID", "chatID", chatID)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		intChatID, err := strconv.ParseInt(chatID, 10, 64)
		if err != nil {
			h.log.Error("handler-GetLinks: could not parse int chatID", "error", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		links, err := h.useCase.GetLinks(ctx, intChatID)
		if err != nil {
			h.log.Error("handler-GetLinks: could not get links", "error", err)

			return
		}

		h.log.Info("handler-GetLinks: Successfully got links")
		render.JSON(w, r, links)
	}
}
