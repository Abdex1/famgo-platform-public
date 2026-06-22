package auth

// Claims represents verified JWT identity used across FamGo services.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
