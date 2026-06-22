package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/Abdex1/FamGo-platform/services/user-service/internal/domain/entities"
)

// UserProfileRepository handles user profile database operations
type UserProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository creates a new user profile repository
func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{db: db}
}

// Create inserts a new user profile
func (r *UserProfileRepository) Create(ctx context.Context, profile *entities.UserProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}
	return nil
}

// GetByUserID retrieves user profile by user ID
func (r *UserProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserProfile, error) {
	var profile entities.UserProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user profile not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return &profile, nil
}

// Update updates an existing user profile
func (r *UserProfileRepository) Update(ctx context.Context, profile *entities.UserProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}
	return nil
}

// Delete deletes a user profile
func (r *UserProfileRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entities.UserProfile{}).Error; err != nil {
		return fmt.Errorf("failed to delete user profile: %w", err)
	}
	return nil
}

// UserPreferenceRepository handles user preference database operations
type UserPreferenceRepository struct {
	db *gorm.DB
}

// NewUserPreferenceRepository creates a new user preference repository
func NewUserPreferenceRepository(db *gorm.DB) *UserPreferenceRepository {
	return &UserPreferenceRepository{db: db}
}

// Create inserts user preferences
func (r *UserPreferenceRepository) Create(ctx context.Context, pref *entities.UserPreference) error {
	if err := r.db.WithContext(ctx).Create(pref).Error; err != nil {
		return fmt.Errorf("failed to create user preference: %w", err)
	}
	return nil
}

// GetByUserID retrieves user preferences
func (r *UserPreferenceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.UserPreference, error) {
	var pref entities.UserPreference
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user preference not found")
		}
		return nil, fmt.Errorf("failed to get user preference: %w", err)
	}
	return &pref, nil
}

// Update updates user preferences
func (r *UserPreferenceRepository) Update(ctx context.Context, pref *entities.UserPreference) error {
	if err := r.db.WithContext(ctx).Save(pref).Error; err != nil {
		return fmt.Errorf("failed to update user preference: %w", err)
	}
	return nil
}

// UserRatingRepository handles user rating database operations
type UserRatingRepository struct {
	db *gorm.DB
}

// NewUserRatingRepository creates a new user rating repository
func NewUserRatingRepository(db *gorm.DB) *UserRatingRepository {
	return &UserRatingRepository{db: db}
}

// Create inserts a new rating
func (r *UserRatingRepository) Create(ctx context.Context, rating *entities.UserRating) error {
	if err := r.db.WithContext(ctx).Create(rating).Error; err != nil {
		return fmt.Errorf("failed to create rating: %w", err)
	}
	return nil
}

// GetByUserID retrieves all ratings for a user
func (r *UserRatingRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entities.UserRating, error) {
	var ratings []entities.UserRating
	if err := r.db.WithContext(ctx).
		Where("rated_user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ratings).Error; err != nil {
		return nil, fmt.Errorf("failed to get ratings: %w", err)
	}
	return ratings, nil
}

// GetAverageRating calculates average rating for a user
func (r *UserRatingRepository) GetAverageRating(ctx context.Context, userID uuid.UUID) (float64, int64, error) {
	var result struct {
		Average float64
		Count   int64
	}
	
	if err := r.db.WithContext(ctx).
		Model(&entities.UserRating{}).
		Where("rated_user_id = ?", userID).
		Select("AVG(rating) as average, COUNT(*) as count").
		Scan(&result).Error; err != nil {
		return 0, 0, fmt.Errorf("failed to get average rating: %w", err)
	}
	
	return result.Average, result.Count, nil
}

// UserNotificationRepository handles user notification database operations
type UserNotificationRepository struct {
	db *gorm.DB
}

// NewUserNotificationRepository creates a new user notification repository
func NewUserNotificationRepository(db *gorm.DB) *UserNotificationRepository {
	return &UserNotificationRepository{db: db}
}

// Create inserts a new notification
func (r *UserNotificationRepository) Create(ctx context.Context, notif *entities.UserNotification) error {
	if err := r.db.WithContext(ctx).Create(notif).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

// GetByUserID retrieves user notifications
func (r *UserNotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entities.UserNotification, error) {
	var notifs []entities.UserNotification
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifs).Error; err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	return notifs, nil
}

// MarkAsRead marks notification as read
func (r *UserNotificationRepository) MarkAsRead(ctx context.Context, notifID uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&entities.UserNotification{}).
		Where("id = ?", notifID).
		Updates(map[string]interface{}{
			"read": true,
			"read_at": now,
		}).Error; err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

// Delete soft deletes a notification
func (r *UserNotificationRepository) Delete(ctx context.Context, notifID uuid.UUID) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&entities.UserNotification{}).
		Where("id = ?", notifID).
		Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// GetUnreadCount gets count of unread notifications
func (r *UserNotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entities.UserNotification{}).
		Where("user_id = ? AND read = false AND deleted_at IS NULL", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}
