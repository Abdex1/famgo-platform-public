package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"famgo/driver-service/internal/domain/entities"
)

// DriverRepository handles driver database operations
type DriverRepository struct {
	db *gorm.DB
}

// NewDriverRepository creates a new driver repository
func NewDriverRepository(db *gorm.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

// Create inserts a new driver
func (r *DriverRepository) Create(ctx context.Context, driver *entities.DriverProfile) error {
	if err := r.db.WithContext(ctx).Create(driver).Error; err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}
	return nil
}

// GetByID retrieves driver by ID
func (r *DriverRepository) GetByID(ctx context.Context, driverID uuid.UUID) (*entities.DriverProfile, error) {
	var driver entities.DriverProfile
	if err := r.db.WithContext(ctx).Where("id = ?", driverID).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("driver not found")
		}
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}
	return &driver, nil
}

// GetByUserID retrieves driver by user ID
func (r *DriverRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entities.DriverProfile, error) {
	var driver entities.DriverProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("driver not found")
		}
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}
	return &driver, nil
}

// Update updates a driver
func (r *DriverRepository) Update(ctx context.Context, driver *entities.DriverProfile) error {
	if err := r.db.WithContext(ctx).Save(driver).Error; err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}
	return nil
}

// UpdateStatus updates driver status
func (r *DriverRepository) UpdateStatus(ctx context.Context, driverID uuid.UUID, status entities.DriverStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == entities.DriverOnline {
		updates["last_online_at"] = time.Now()
	}
	if err := r.db.WithContext(ctx).Model(&entities.DriverProfile{}).Where("id = ?", driverID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update driver status: %w", err)
	}
	return nil
}

// GetNearbyDrivers retrieves drivers within radius
func (r *DriverRepository) GetNearbyDrivers(ctx context.Context, latitude, longitude float64, radiusKm int) ([]entities.DriverProfile, error) {
	var drivers []entities.DriverProfile
	// TODO: Use PostGIS to query nearby drivers
	return drivers, nil
}

// VehicleRepository handles vehicle database operations
type VehicleRepository struct {
	db *gorm.DB
}

// NewVehicleRepository creates a new vehicle repository
func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// Create inserts a new vehicle
func (r *VehicleRepository) Create(ctx context.Context, vehicle *entities.Vehicle) error {
	if err := r.db.WithContext(ctx).Create(vehicle).Error; err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}
	return nil
}

// GetByID retrieves vehicle by ID
func (r *VehicleRepository) GetByID(ctx context.Context, vehicleID uuid.UUID) (*entities.Vehicle, error) {
	var vehicle entities.Vehicle
	if err := r.db.WithContext(ctx).Where("id = ?", vehicleID).First(&vehicle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}
	return &vehicle, nil
}

// GetByDriverID retrieves vehicles for a driver
func (r *VehicleRepository) GetByDriverID(ctx context.Context, driverID uuid.UUID) ([]entities.Vehicle, error) {
	var vehicles []entities.Vehicle
	if err := r.db.WithContext(ctx).Where("driver_id = ? AND is_active = true", driverID).Find(&vehicles).Error; err != nil {
		return nil, fmt.Errorf("failed to get vehicles: %w", err)
	}
	return vehicles, nil
}

// Update updates a vehicle
func (r *VehicleRepository) Update(ctx context.Context, vehicle *entities.Vehicle) error {
	if err := r.db.WithContext(ctx).Save(vehicle).Error; err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	return nil
}

// Delete soft deletes a vehicle
func (r *VehicleRepository) Delete(ctx context.Context, vehicleID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&entities.Vehicle{}).Where("id = ?", vehicleID).Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}
	return nil
}

// DriverLocationRepository handles driver location operations
type DriverLocationRepository struct {
	db *gorm.DB
}

// NewDriverLocationRepository creates a new driver location repository
func NewDriverLocationRepository(db *gorm.DB) *DriverLocationRepository {
	return &DriverLocationRepository{db: db}
}

// UpdateLocation updates or creates driver location
func (r *DriverLocationRepository) UpdateLocation(ctx context.Context, location *entities.DriverLocation) error {
	if err := r.db.WithContext(ctx).Save(location).Error; err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}
	return nil
}

// GetByDriverID retrieves driver location
func (r *DriverLocationRepository) GetByDriverID(ctx context.Context, driverID uuid.UUID) (*entities.DriverLocation, error) {
	var location entities.DriverLocation
	if err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).First(&location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("location not found")
		}
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	return &location, nil
}

// DriverEarningsRepository handles earnings data
type DriverEarningsRepository struct {
	db *gorm.DB
}

// NewDriverEarningsRepository creates a new driver earnings repository
func NewDriverEarningsRepository(db *gorm.DB) *DriverEarningsRepository {
	return &DriverEarningsRepository{db: db}
}

// GetByDriverAndPeriod retrieves earnings for period
func (r *DriverEarningsRepository) GetByDriverAndPeriod(ctx context.Context, driverID uuid.UUID, period string, start, end time.Time) (*entities.DriverEarnings, error) {
	var earnings entities.DriverEarnings
	if err := r.db.WithContext(ctx).
		Where("driver_id = ? AND period = ? AND period_start = ? AND period_end = ?", driverID, period, start, end).
		First(&earnings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("earnings not found")
		}
		return nil, fmt.Errorf("failed to get earnings: %w", err)
	}
	return &earnings, nil
}

// GetTotalEarnings retrieves total earnings for period
func (r *DriverEarningsRepository) GetTotalEarnings(ctx context.Context, driverID uuid.UUID, days int) (interface{}, error) {
	type Result struct {
		TotalEarnings  float64
		CompletedRides int
		CancelledRides int
	}

	var result Result
	startDate := time.Now().AddDate(0, 0, -days)

	if err := r.db.WithContext(ctx).
		Model(&entities.DriverEarnings{}).
		Where("driver_id = ? AND period_end >= ?", driverID, startDate).
		Select("SUM(net_earnings) as total_earnings, SUM(rides_completed) as completed_rides, SUM(rides_cancelled) as cancelled_rides").
		Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get total earnings: %w", err)
	}

	return result, nil
}
