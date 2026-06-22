/*
# PHASE 17 — ENTERPRISE AUTH ROUTES

internal/interfaces/rest/routes/routes.go
*/
package routes

import (
    "github.com/go-chi/chi/v5"

    "famgo/auth-service/internal/interfaces/rest/handlers"
)

func RegisterRoutes(
    r chi.Router,
) {

    authHandler := handlers.NewAuthHandler()

    r.Get("/health", handlers.Health)

    r.Route("/api/v1/auth", func(r chi.Router) {

        r.Post("/register", authHandler.Register)

    })
}
