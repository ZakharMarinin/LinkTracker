package router

import (
	"context"
	"log/slog"
	"scrapper/internal/http/handlers"
	"scrapper/internal/http/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(ctx context.Context, router *chi.Mux, http *handlers.HTTP, log *slog.Logger) {
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

	router.Post("/tg-chat/{id}", http.CreateChat(ctx))
	router.Delete("/tg-chat/{id}", http.DeleteChat(ctx))
	router.Get("/links/{id}", http.GetLinks(ctx))
	router.Post("/links", http.AddLink(ctx))
	router.Delete("/links", http.DeleteLink(ctx))
}
