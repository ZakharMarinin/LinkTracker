package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"scrapper/internal/domain"

	"github.com/go-chi/render"
)

func (h *HTTP) AddLink(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Error("error reading body", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var link *domain.AddLinkRequest

		err = json.Unmarshal(body, &link)
		if err != nil {
			h.log.Error("error unmarshalling body", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.useCase.AddLink(ctx, link.ChatID, link.URL, link.Desc, link.Tags)
		if err != nil {
			h.log.Error("error adding link", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.log.Info("successfully added link")
		w.WriteHeader(http.StatusOK)

		render.JSON(w, r, "successfully added link")
	}
}
