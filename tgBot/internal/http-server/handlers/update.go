package handlers

import (
	"context"
	"errors"
	"linktracker/internal/domain"
	"linktracker/internal/domain/api/response"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	response.Response
}

type HTTP struct {
	TGService URLUpdate
	log       *slog.Logger
}

type URLUpdate interface {
	Updates(update *domain.UpdatedLink) error
}

func NewURLUpdate(update URLUpdate, log *slog.Logger) *HTTP {
	return &HTTP{update, log}
}

func (h *HTTP) SendUpdates(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.SendUpdates"
		h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req domain.UpdatedLink

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Error("SendUpdates: Error unmarshalling request", "error", err)
			render.JSON(w, r, response.Error("SendUpdates: Error unmarshalling request"))
			return
		}

		h.log.Info("SendUpdates: request body decoded", "req", req)

		err = validator.New().Struct(req)
		if err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			h.log.Error("SendUpdates: Error validating request", "error", err)
			render.JSON(w, r, response.ValidatorError(validateErr))
			return
		}

		err = h.TGService.Updates(&req)
		if err != nil {
			h.log.Error("SendUpdates: Error updating links", "error", err)
			return
		}
	}
}
