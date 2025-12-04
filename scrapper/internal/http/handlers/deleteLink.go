package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"scrapper/internal/domain"

	"github.com/go-chi/render"
)

func (h *HTTP) DeleteLink(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Error("error reading body", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var link *domain.DeleteLinkRequest

		err = json.Unmarshal(body, &link)
		if err != nil {
			h.log.Error("error unmarshalling body", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.useCase.DeleteLink(ctx, link.ChatID, link.Alias)
		if err != nil {
			h.log.Error("error deleting link", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.log.Info("successfully deleted link")
		w.WriteHeader(http.StatusOK)

		render.JSON(w, r, "successfully deleted link")
	}
}
