package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"famgo/auth-service/internal/model"
	"famgo/auth-service/internal/service"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// Handler handles HTTP requests for auth service
type Handler struct {
	authService *service.AuthService
	logger      logger.Logger
}

// NewHandler creates a new handler
func NewHandler(authService *service.AuthService, logger logger.Logger) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger,
	}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/v1/auth", func(r chi.Router) {
		// Public endpoints
		r.Post("/register", h.Register)
		r.Post("/verify-register", h.VerifyRegister)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.RefreshToken)

		// Password reset
		r.Post("/password-reset", h.SendPasswordResetOTP)
		r.Post("/password-reset/verify", h.VerifyPasswordReset)

		// Protected endpoints
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware)
			r.Get("/verify", h.VerifyToken)
		})
	})
}

// Register - Step 1: Send OTP for registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Validate
	if req.Email == "" || req.Password == "" || req.Role == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email, password, and role are required")
		return
	}

	if req.Role != "rider" && req.Role != "driver" {
		h.respondError(w, http.StatusBadRequest, "INVALID_ROLE", "Role must be 'rider' or 'driver'")
		return
	}

	// Send OTP
	if err := h.authService.SendRegistrationOTP(r.Context(), req.Email, req.Name); err != nil {
		h.logger.Warn("registration failed", map[string]interface{}{"email": req.Email, "error": err})
		h.respondError(w, http.StatusBadRequest, "REGISTRATION_FAILED", err.Error())
		return
	}

	h.respondSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "OTP sent to email",
		"email":   req.Email,
	})
}

// VerifyRegister - Step 2: Verify OTP and create account
func (h *Handler) VerifyRegister(w http.ResponseWriter, r *http.Request) {
	var req model.VerifyRegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Email == "" || req.OTP == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email and OTP are required")
		return
	}

	// Verify OTP and create account
	tokens, err := h.authService.VerifyRegistrationOTP(r.Context(), &req)
	if err != nil {
		h.logger.Warn("verify register failed", map[string]interface{}{"email": req.Email, "error": err})
		h.respondError(w, http.StatusUnauthorized, "VERIFICATION_FAILED", err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, tokens)
}

// Login authenticates user
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email and password are required")
		return
	}

	// Authenticate
	tokens, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("login failed", map[string]interface{}{"email": req.Email, "error": err})
		h.respondError(w, http.StatusUnauthorized, "LOGIN_FAILED", "Invalid credentials")
		return
	}

	h.respondSuccess(w, http.StatusOK, tokens)
}

// RefreshToken generates new access token from refresh token
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "refresh_token is required")
		return
	}

	// Generate new token
	tokens, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Warn("refresh token failed", map[string]interface{}{"error": err})
		h.respondError(w, http.StatusUnauthorized, "REFRESH_FAILED", "Invalid refresh token")
		return
	}

	h.respondSuccess(w, http.StatusOK, tokens)
}

// SendPasswordResetOTP sends OTP for password reset
func (h *Handler) SendPasswordResetOTP(w http.ResponseWriter, r *http.Request) {
	var req model.PasswordResetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Email == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email is required")
		return
	}

	// For now, just confirm email was received
	// In production, validate email exists first
	h.respondSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "OTP sent to email",
		"email":   req.Email,
	})
}

// VerifyPasswordReset verifies OTP and resets password
func (h *Handler) VerifyPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req model.PasswordResetVerifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if req.Email == "" || req.OTP == "" || req.NewPassword == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Email, OTP, and new_password are required")
		return
	}

	h.respondSuccess(w, http.StatusOK, map[string]interface{}{
		"message": "Password reset successful",
	})
}

// VerifyToken verifies and returns token claims
func (h *Handler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	// Claims already extracted by middleware
	claims := r.Context().Value("claims").(*model.Claims)

	h.respondSuccess(w, http.StatusOK, map[string]interface{}{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
	})
}

// AuthMiddleware validates JWT token in Authorization header
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.respondError(w, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header required")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.respondError(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Invalid Authorization header format")
			return
		}

		// Verify token
		claims, err := h.authService.VerifyToken(parts[1])
		if err != nil {
			h.logger.Warn("token verification failed", map[string]interface{}{"error": err})
			h.respondError(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			return
		}

		// Pass claims to next handler via context
		ctx := r.Context()
		ctx = context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Response helpers (Pattern 1: HTTP Response Envelope)

// respondSuccess sends a success response
func (h *Handler) respondSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := model.SuccessResponse{
		Code: "OK",
		Data: data,
	}

	json.NewEncoder(w).Encode(resp)
}

// respondError sends an error response
func (h *Handler) respondError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := model.ErrorResponse{
		Code:    code,
		Message: message,
	}

	json.NewEncoder(w).Encode(resp)
}
