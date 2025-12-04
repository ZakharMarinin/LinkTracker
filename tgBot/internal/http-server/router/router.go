package router

import (
	"context"
	"linktracker/internal/http-server/handlers"
	"linktracker/internal/http-server/middleware/logger"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(router *chi.Mux, HTTP *handlers.HTTP, ctx context.Context, log *slog.Logger) {
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

	router.Post("/updates", HTTP.SendUpdates(ctx))
}
