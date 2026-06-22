/*
# PHASE 18 — MAIN BOOTSTRAP

cmd/api/main.go
*/
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"

    "famgo/auth-service/internal/interfaces/rest/routes"
)

func main() {

    r := chi.NewRouter()

    routes.RegisterRoutes(r)

    http.ListenAndServe(":8080", r)
}
