/*
# PHASE 10 — CREATE LOGIN DTOs

internal/application/dto/auth_dto.go
*/
package dto

type RegisterRequest struct {
    Email       string `json:"email"`
    Password    string `json:"password"`
}

type LoginRequest struct {
    Email       string `json:"email"`
    Password    string `json:"password"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}
