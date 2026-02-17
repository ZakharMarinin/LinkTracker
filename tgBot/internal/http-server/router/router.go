package router

import (
	"linktracker/internal/http-server/handlers"
	"linktracker/internal/http-server/middleware/logger"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router(router *chi.Mux, http *handlers.HTTP, log *slog.Logger) {
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

	router.Post("/updates", http.SendUpdates())
	router.Handle("/metrics", promhttp.Handler())
}
