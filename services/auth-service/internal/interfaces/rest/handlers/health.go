/*
# PHASE 15 — HEALTH CHECKS

internal/interfaces/rest/handlers/health.go
*/
package handlers

import (
    "encoding/json"
    "net/http"
)

func Health(
    w http.ResponseWriter,
    r *http.Request,
) {

    response := map[string]string{
        "status": "ok",
        "service": "github.com/Abdex1/FamGo-platform/services/auth-service",
    }

    w.Header().Set("Content-Type", "application/json")

    json.NewEncoder(w).Encode(response)
}
