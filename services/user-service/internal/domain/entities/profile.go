package entities

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// UserProfile represents extended user profile information
type UserProfile struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex"`
	Bio             string         `gorm:"type:text"`
	ProfileImageURL string
	CoverImageURL   string
	Address         string
	City            string
	State           string
	ZipCode         string
	Country         string
	PhoneVerified   bool           `gorm:"default:false"`
	EmailVerified   bool           `gorm:"default:false"`
	PhoneVerifiedAt *time.Time
	EmailVerifiedAt *time.Time
	Languages       datatypes.JSONSlice `gorm:"type:jsonb"` // ["Amharic", "English", "Tigrinya"]
	Preferences     datatypes.JSONMap   `gorm:"type:jsonb;default:'{}'"`
	EmergencyContact *json.RawMessage   `gorm:"type:jsonb"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
}

// UserRating represents user ratings
type UserRating struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	RatedByID  uuid.UUID `gorm:"type:uuid;not null"`
	RatedUserID uuid.UUID `gorm:"type:uuid;not null"`
	RideID     uuid.UUID `gorm:"type:uuid"`
	Rating     int       `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Review     string    `gorm:"type:text"`
	Categories datatypes.JSONMap `gorm:"type:jsonb"` // {cleanliness: 5, comfort: 4}
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// UserPreference represents user preferences
type UserPreference struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Language          string    `gorm:"default:'en'"`
	Currency          string    `gorm:"default:'ETB'"`
	ThemMode          string    `gorm:"default:'light'"` // light, dark, auto
	NotificationsEnabled bool   `gorm:"default:true"`
	EmailNotifications bool    `gorm:"default:true"`
	SMSNotifications  bool    `gorm:"default:true"`
	PushNotifications bool    `gorm:"default:true"`
	ShareRideData    bool    `gorm:"default:false"`
	EmergencySharing bool    `gorm:"default:false"`
	FemaleDriverPreference bool `gorm:"default:false"`
	PoolingPreference string  `gorm:"default:'standard'"` // standard, pool, both
	RideTypePreference string `gorm:"default:'all'"` // standard, comfort, xl, all
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}

// UserActivity represents user activity tracking
type UserActivity struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	ActivityType string    // ride_requested, ride_completed, payment_made, etc.
	Reference    datatypes.JSONMap `gorm:"type:jsonb"`
	IPAddress    string
	UserAgent    string
	Timestamp    time.Time `gorm:"autoCreateTime;index"`
}

// UserNotification represents user notifications
type UserNotification struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Type       string    // ride_update, payment, promotion, system
	Title      string
	Body       string
	ActionURL  string
	ImageURL   string
	Data       datatypes.JSONMap `gorm:"type:jsonb"`
	Read       bool              `gorm:"default:false"`
	ReadAt     *time.Time
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	DeletedAt  *time.Time `gorm:"index"`
}

// TableName specifies table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profiles"
}

// TableName specifies table name for UserRating
func (UserRating) TableName() string {
	return "user_ratings"
}

// TableName specifies table name for UserPreference
func (UserPreference) TableName() string {
	return "user_preferences"
}

// TableName specifies table name for UserActivity
func (UserActivity) TableName() string {
	return "user_activity"
}

// TableName specifies table name for UserNotification
func (UserNotification) TableName() string {
	return "user_notifications"
}

// Scan implements sql.Scanner interface
func (up *UserProfile) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &up)
}

// Value implements driver.Valuer interface
func (up UserProfile) Value() (driver.Value, error) {
	return json.Marshal(up)
}
