package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// NotificationChannel represents notification delivery channel
type NotificationChannel string

const (
	ChannelSMS   NotificationChannel = "sms"
	ChannelPush  NotificationChannel = "push"
	ChannelEmail NotificationChannel = "email"
	ChannelInApp NotificationChannel = "in_app"
)

// NotificationStatus represents notification delivery status
type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusQueued    NotificationStatus = "queued"
	StatusSent      NotificationStatus = "sent"
	StatusDelivered NotificationStatus = "delivered"
	StatusFailed    NotificationStatus = "failed"
)

// Notification represents a notification record
type Notification struct {
	ID              uuid.UUID           `gorm:"type:uuid;primaryKey"`
	UserID          uuid.UUID           `gorm:"type:uuid;not null;index"`
	Channel         NotificationChannel `gorm:"type:notification_channel"`
	RecipientPhone  string              // For SMS
	RecipientEmail  string              // For email
	Title           string
	Body            string
	ActionURL       string
	ImageURL        string
	Data            datatypes.JSONMap   `gorm:"type:jsonb"` // Additional data
	Status          NotificationStatus  `gorm:"default:'pending'"`
	Provider        string              // twilio, firebase, sendgrid
	ProviderID      string              // Reference from provider
	SentAt          *time.Time
	DeliveredAt     *time.Time
	ReadAt          *time.Time
	ErrorMessage    string
	RetryCount      int                 `gorm:"default:0"`
	MaxRetries      int                 `gorm:"default:3"`
	NextRetryAt     *time.Time
	CreatedAt       time.Time           `gorm:"autoCreateTime"`
	UpdatedAt       time.Time           `gorm:"autoUpdateTime"`
}

// NotificationTemplate represents notification template for campaigns
type NotificationTemplate struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey"`
	Name        string          `gorm:"not null"`
	Description string
	Channel     NotificationChannel
	Title       string
	BodyTemplate string // With placeholders: {user_name}, {ride_id}, etc.
	ActionURL   string
	Variables   datatypes.JSONSlice `gorm:"type:jsonb"` // List of variable names
	IsActive    bool            `gorm:"default:true"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
}

// NotificationPreference represents user notification preferences
type NotificationPreference struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID               uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	SMSEnabled           bool      `gorm:"default:true"`
	PushEnabled          bool      `gorm:"default:true"`
	EmailEnabled         bool      `gorm:"default:false"`
	RideUpdatesSMS       bool      `gorm:"default:true"`
	RideUpdatesPush      bool      `gorm:"default:true"`
	PromoNotifications   bool      `gorm:"default:false"`
	SafetyAlerts         bool      `gorm:"default:true"`
	PaymentNotifications bool      `gorm:"default:true"`
	EarningsUpdates      bool      `gorm:"default:true"` // For drivers
	QuietHoursStart      string    // HH:MM format
	QuietHoursEnd        string
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
}

// NotificationLog represents delivery log
type NotificationLog struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	NotificationID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Channel         NotificationChannel
	Status          NotificationStatus
	ResponseCode    string
	ResponseMessage string
	DeliveryTime    int       // milliseconds
	Timestamp       time.Time `gorm:"autoCreateTime"`
}

// TableName specifies table name
func (Notification) TableName() string {
	return "notifications"
}

// TableName specifies table name
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

// TableName specifies table name
func (NotificationPreference) TableName() string {
	return "notification_preferences"
}

// TableName specifies table name
func (NotificationLog) TableName() string {
	return "notification_logs"
}
