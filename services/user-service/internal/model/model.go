package model

import "time"

// UserProfile represents a user profile
type UserProfile struct {
	ID               string    `db:"id"`
	AuthID           string    `db:"auth_id"`
	FirstName        string    `db:"first_name"`
	LastName         string    `db:"last_name"`
	ProfilePictureURL string   `db:"profile_picture_url"`
	EmailVerified    bool      `db:"email_verified"`
	PhoneVerified    bool      `db:"phone_verified"`
	Rating           float64   `db:"rating"`
	TotalRides       int       `db:"total_rides"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// UserPreferences represents user preferences
type UserPreferences struct {
	ID               string    `db:"id"`
	UserID           string    `db:"user_id"`
	NotificationEmail bool     `db:"notification_email"`
	NotificationSMS  bool      `db:"notification_sms"`
	Language         string    `db:"language"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// UserAddress represents saved addresses
type UserAddress struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	Type         string    `db:"type"` // "home", "work"
	AddressLine1 string    `db:"address_line_1"`
	City         string    `db:"city"`
	Lat          float64   `db:"lat"`
	Lng          float64   `db:"lng"`
	CreatedAt    time.Time `db:"created_at"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

// UpdatePreferencesRequest represents preferences update request
type UpdatePreferencesRequest struct {
	NotificationEmail bool   `json:"notification_email"`
	NotificationSMS  bool    `json:"notification_sms"`
	Language         string  `json:"language"`
}

// AddressRequest represents address creation/update request
type AddressRequest struct {
	Type         string  `json:"type"`
	AddressLine1 string  `json:"address_line_1"`
	City         string  `json:"city"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
}
