
//shared/contracts/events/payment/v1/payment_completed.go

package v1

import "time"

type PaymentCompleted struct {
    PaymentID   string    `json:"payment_id"`
    TripID      string    `json:"trip_id"`
    RiderID     string    `json:"rider_id"`
    DriverID    string    `json:"driver_id"`
    Amount      float64   `json:"amount"`
    Currency    string    `json:"currency"`
    Status      string    `json:"status"`
    CompletedAt time.Time `json:"completed_at"`
}
