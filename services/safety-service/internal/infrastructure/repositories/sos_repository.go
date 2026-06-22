// services/safety-service/internal/infrastructure/repositories/sos_repository.go
package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Abdex1/FamGo-platform/services/safety-service/internal/domain/entities"
)

type SOSRepository struct {
	pool *pgxpool.Pool
}

func NewSOSRepository(pool *pgxpool.Pool) *SOSRepository {
	return &SOSRepository{pool: pool}
}

func (r *SOSRepository) Create(ctx context.Context, incident *entities.SOSIncident) error {
	query := `
		INSERT INTO sos_incidents 
		(id, ride_id, user_id, user_type, latitude, longitude, status, 
		 incident_type, description, escalation_level, initiated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.pool.Exec(ctx, query,
		incident.ID, incident.RideID, incident.UserID, incident.UserType,
		incident.Latitude, incident.Longitude, string(incident.Status),
		incident.IncidentType, incident.Description, string(incident.EscalationLevel),
		incident.InitiatedAt, incident.CreatedAt, incident.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create SOS incident: %w", err)
	}
	return nil
}

func (r *SOSRepository) GetByID(ctx context.Context, id string) (*entities.SOSIncident, error) {
	query := `
		SELECT id, ride_id, user_id, user_type, latitude, longitude, status,
		       incident_type, description, escalation_level, notifications_sent,
		       police_notified, police_notified_at, ambulance_notified, ambulance_notified_at,
		       platform_notified, platform_notified_at, resolved_at, resolved_by,
		       resolution_notes, initiated_at, last_escalated_at, created_at, updated_at
		FROM sos_incidents WHERE id = $1 AND deleted_at IS NULL
	`
	row := r.pool.QueryRow(ctx, query, id)

	var incident entities.SOSIncident
	err := row.Scan(
		&incident.ID, &incident.RideID, &incident.UserID, &incident.UserType,
		&incident.Latitude, &incident.Longitude, &incident.Status,
		&incident.IncidentType, &incident.Description, &incident.EscalationLevel,
		&incident.NotificationsSent, &incident.PoliceNotified, &incident.PoliceNotifiedAt,
		&incident.AmbulanceNotified, &incident.AmbulanceNotifiedAt,
		&incident.PlatformNotified, &incident.PlatformNotifiedAt, &incident.ResolvedAt,
		&incident.ResolvedBy, &incident.ResolutionNotes, &incident.InitiatedAt,
		&incident.LastEscalatedAt, &incident.CreatedAt, &incident.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("incident not found")
		}
		return nil, fmt.Errorf("failed to query incident: %w", err)
	}
	return &incident, nil
}

func (r *SOSRepository) Update(ctx context.Context, incident *entities.SOSIncident) error {
	query := `
		UPDATE sos_incidents SET
			status = $1, escalation_level = $2, notifications_sent = $3,
			police_notified = $4, police_notified_at = $5, ambulance_notified = $6,
			ambulance_notified_at = $7, platform_notified = $8, platform_notified_at = $9,
			resolved_at = $10, resolved_by = $11, resolution_notes = $12,
			last_escalated_at = $13, updated_at = $14
		WHERE id = $15 AND deleted_at IS NULL
	`
	result, err := r.pool.Exec(ctx, query,
		string(incident.Status), string(incident.EscalationLevel), incident.NotificationsSent,
		incident.PoliceNotified, incident.PoliceNotifiedAt, incident.AmbulanceNotified,
		incident.AmbulanceNotifiedAt, incident.PlatformNotified, incident.PlatformNotifiedAt,
		incident.ResolvedAt, incident.ResolvedBy, incident.ResolutionNotes,
		incident.LastEscalatedAt, incident.UpdatedAt, incident.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update incident: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("incident not found")
	}
	return nil
}

func (r *SOSRepository) GetActiveIncidents(ctx context.Context) ([]*entities.SOSIncident, error) {
	query := `
		SELECT id, ride_id, user_id, user_type, latitude, longitude, status,
		       incident_type, description, escalation_level, notifications_sent,
		       police_notified, police_notified_at, ambulance_notified, ambulance_notified_at,
		       platform_notified, platform_notified_at, resolved_at, resolved_by,
		       resolution_notes, initiated_at, last_escalated_at, created_at, updated_at
		FROM sos_incidents WHERE status IN ('initiated', 'active', 'escalated')
		  AND deleted_at IS NULL
		ORDER BY initiated_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query incidents: %w", err)
	}
	defer rows.Close()

	var incidents []*entities.SOSIncident

	for rows.Next() {
		var incident entities.SOSIncident
		err := rows.Scan(
			&incident.ID, &incident.RideID, &incident.UserID, &incident.UserType,
			&incident.Latitude, &incident.Longitude, &incident.Status,
			&incident.IncidentType, &incident.Description, &incident.EscalationLevel,
			&incident.NotificationsSent, &incident.PoliceNotified, &incident.PoliceNotifiedAt,
			&incident.AmbulanceNotified, &incident.AmbulanceNotifiedAt,
			&incident.PlatformNotified, &incident.PlatformNotifiedAt, &incident.ResolvedAt,
			&incident.ResolvedBy, &incident.ResolutionNotes, &incident.InitiatedAt,
			&incident.LastEscalatedAt, &incident.CreatedAt, &incident.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}
		incidents = append(incidents, &incident)
	}

	return incidents, rows.Err()
}
