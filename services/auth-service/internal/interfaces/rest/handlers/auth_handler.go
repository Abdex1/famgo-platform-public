/*
# PHASE 11 — CREATE REST HANDLERS

internal/interfaces/rest/handlers/auth_handler.go

go get golang.org/x/crypto/bcrypt

*/
package handlers

import (
    "encoding/json"
    "net/http"

    "golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{}
}

func (h *AuthHandler) Register(
    w http.ResponseWriter,
    r *http.Request,
) {

    type req struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    var body req

    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    hash, err := bcrypt.GenerateFromPassword(
        []byte(body.Password),
        bcrypt.DefaultCost,
    )

    if err != nil {
        http.Error(w, "failed", http.StatusInternalServerError)
        return
    }

    _ = hash

    w.WriteHeader(http.StatusCreated)
}
