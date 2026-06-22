// Location: services/payment-service/internal/domain/payment_intent.go

type PaymentIntent struct {
    ID                  string    `db:"id"`
    RideID              string    `db:"ride_id"`
    UserID              string    `db:"user_id"`
    Amount              float64   `db:"amount"`
    Currency            string    `db:"currency"`
    Status              string    `db:"status"` // pending, processing, succeeded, failed
    Provider            string    `db:"provider"` // telebirr, cbe, chapa
    ProviderIntentID    string    `db:"provider_intent_id"`
    
    CreatedAt           time.Time `db:"created_at"`
    ExpiresAt           time.Time `db:"expires_at"` // 15 minutes
    SucceededAt         *time.Time `db:"succeeded_at"`
    FailedAt            *time.Time `db:"failed_at"`
    FailureReason       *string   `db:"failure_reason"`
    
    RetryCount          int       `db:"retry_count"`
    MaxRetries          int       `db:"max_retries"` // 3
}

func (s *PaymentService) CreatePaymentIntent(
    ctx context.Context,
    ride *Ride,
    paymentMethod string,
) (*PaymentIntent, error) {
    
    intent := &PaymentIntent{
        ID:          uuid.New().String(),
        RideID:      ride.ID,
        UserID:      ride.UserID,
        Amount:      ride.FareAmount,
        Currency:    "ETB",
        Status:      "pending",
        Provider:    s.selectProvider(paymentMethod),
        CreatedAt:   time.Now(),
        ExpiresAt:   time.Now().Add(15 * time.Minute),
        RetryCount:  0,
        MaxRetries:  3,
    }
    
    // Store intent
    s.repo.InsertIntent(ctx, intent)
    
    // Publish event
    s.eventBus.Publish(ctx, "payment.intent_created", intent)
    
    return intent, nil
}

// Result: ✅ Payment intent created