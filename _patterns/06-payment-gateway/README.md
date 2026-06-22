# 💳 PAYMENT GATEWAY PATTERNS
## Extracted from uber-master, Adapted for FamGo

**Status:** Pattern 6/8

---

## Gateway Abstraction

```go
type PaymentGateway interface {
    CreateOrder(ctx context.Context, amount int64, currency string) (*Order, error)
    VerifySignature(signature, payload string) (bool, error)
    Refund(ctx context.Context, transactionID string) error
}

// Razorpay Implementation
type RazorpayGateway struct {
    keyID     string
    keySecret string
}

// For FamGo, implement:
// - TelebeirGateway
// - CBEBirrGateway
// - ChapaGateway
// - CashGateway
```

## Provider Factory

```go
func NewPaymentGateway(provider string, config Config) (PaymentGateway, error) {
    switch provider {
    case "razorpay":
        return NewRazorpayGateway(config.RazorpayKey, config.RazorpaySecret)
    case "telebirr":
        return NewTelebeirGateway(config.TelebeirKey)
    case "chapa":
        return NewChapaGateway(config.ChapaKey)
    case "cash":
        return NewCashGateway()
    default:
        return nil, fmt.Errorf("unknown provider: %s", provider)
    }
}
```

## Webhook Handler

```go
func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
    signature := r.Header.Get("X-Razorpay-Signature")
    
    body, _ := io.ReadAll(r.Body)
    
    if valid, _ := h.gateway.VerifySignature(signature, string(body)); !valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    
    var payload map[string]interface{}
    json.Unmarshal(body, &payload)
    
    h.service.HandlePaymentWebhook(r.Context(), payload)
    
    w.WriteHeader(http.StatusOK)
}
```

**Pattern 6 Status:** READY FOR USE

---
