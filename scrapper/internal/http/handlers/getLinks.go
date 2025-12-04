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

		links, err := h.useCase.GetLinks(ctx, intChatID)
		if err != nil {
			h.log.Error("GetLinks: could not get links" + err.Error())
			return
		}

		h.log.Info("Successfully got links")
		render.JSON(w, r, links)
	}
}
