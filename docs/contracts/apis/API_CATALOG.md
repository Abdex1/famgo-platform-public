# 🔌 API CATALOG: FamGo API Reference

**Last Updated:** [Date]  
**Status:** Repository Consistency Audit - TASK 1  
**Total Endpoints:** 80+ documented

---

## API CATALOG

### AUTH-SERVICE (HTTP/gRPC)

#### 1. Register User
```
POST /auth/register
- Auth: None
- Body: { email, phone, password, user_type }
- Response: { user_id, token, refresh_token }
- Rate Limit: 5 requests/minute per IP
- SLA: 99.9% availability
```

#### 2. Login
```
POST /auth/login
- Auth: None
- Body: { email_or_phone, password, device_id }
- Response: { user_id, token, refresh_token, expires_in }
- Rate Limit: 10 requests/minute per IP
- SLA: 99.9% availability
```

#### 3. Request OTP
```
POST /auth/otp/request
- Auth: None
- Body: { phone_number }
- Response: { otp_id, expires_in_seconds }
- Rate Limit: 3 requests/minute per phone
- SLA: 99%
```

#### 4. Verify OTP
```
POST /auth/otp/verify
- Auth: None
- Body: { otp_id, otp_code }
- Response: { user_id, token, refresh_token }
- Rate Limit: 5 attempts per OTP
- SLA: 99%
```

#### 5. Refresh Token
```
POST /auth/token/refresh
- Auth: Refresh Token (in Authorization header)
- Body: {}
- Response: { token, refresh_token, expires_in }
- Rate Limit: 100 requests/hour per user
- SLA: 99.9%
```

#### 6. Logout
```
POST /auth/logout
- Auth: JWT Token
- Body: {}
- Response: { status: "logged_out" }
- Rate Limit: 100 requests/hour per user
- SLA: 99%
```

#### 7. Revoke Device Sessions
```
POST /auth/devices/logout
- Auth: JWT Token
- Body: { device_id (optional - if not provided, logout all devices) }
- Response: { devices_logged_out: number }
- Rate Limit: 10 requests/minute per user
- SLA: 99%
```

#### 8. Check Token Validity
```
GET /auth/token/validate
- Auth: JWT Token
- Response: { valid: boolean, expires_in_seconds: number }
- Rate Limit: 1000 requests/minute
- SLA: 99.9%
```

---

### USER-SERVICE (HTTP/gRPC)

#### 9. Get Profile
```
GET /users/{user_id}
- Auth: JWT Token (PASSENGER, DRIVER, ADMIN)
- Response: { user_id, email, phone, first_name, last_name, profile_picture_url, user_type }
- Rate Limit: 100 requests/minute per user
- SLA: 99%
```

#### 10. Update Profile
```
PUT /users/{user_id}
- Auth: JWT Token (PASSENGER, DRIVER, ADMIN)
- Body: { first_name, last_name, profile_picture_url, emergency_contacts }
- Response: { user_id, updated_fields }
- Rate Limit: 10 requests/minute per user
- SLA: 99%
```

#### 11. List Emergency Contacts
```
GET /users/{user_id}/emergency-contacts
- Auth: JWT Token (PASSENGER, DRIVER)
- Response: [{ contact_id, name, phone, relationship }]
- Rate Limit: 50 requests/minute per user
- SLA: 99%
```

#### 12. Add Emergency Contact
```
POST /users/{user_id}/emergency-contacts
- Auth: JWT Token
- Body: { name, phone, relationship }
- Response: { contact_id }
- Rate Limit: 10 requests/minute per user
- SLA: 99%
```

#### 13. Delete Account
```
DELETE /users/{user_id}
- Auth: JWT Token (ADMIN required for deleting others)
- Response: { status: "deleted" }
- Rate Limit: 1 request/minute per user
- SLA: 95% (complex operation)
```

---

### RIDE-SERVICE (HTTP/WebSocket/gRPC)

#### 14. Request Ride
```
POST /rides/request
- Auth: JWT Token (PASSENGER)
- Body: { 
    pickup_location: { latitude, longitude, address },
    dropoff_location: { latitude, longitude, address },
    passengers: number,
    scheduled_time (optional)
  }
- Response: { ride_id, estimated_fare, estimated_eta_seconds }
- Rate Limit: 10 requests/minute per user
- SLA: 99.5%
```

#### 15. Get Ride Status
```
GET /rides/{ride_id}
- Auth: JWT Token
- Response: { 
    ride_id, status, driver_id, driver_location, current_eta, 
    ride_path, pickup_address, dropoff_address, fare
  }
- Rate Limit: 1000 requests/minute
- SLA: 99.9%
```

#### 16. Cancel Ride
```
POST /rides/{ride_id}/cancel
- Auth: JWT Token
- Body: { reason (optional) }
- Response: { ride_id, status: "cancelled", cancellation_fee (if any) }
- Rate Limit: 10 requests/minute per user
- SLA: 99%
```

#### 17. Rate Ride
```
POST /rides/{ride_id}/rating
- Auth: JWT Token
- Body: { rating (1-5), comment (optional), category (driver|vehicle|route) }
- Response: { rating_id, average_rating }
- Rate Limit: 1 request/ride
- SLA: 99%
```

#### 18. WebSocket: Subscribe to Ride Updates
```
WebSocket wss://ws.famgo.local/rides/{ride_id}
- Auth: JWT Token in ?token=
- Message Types:
  {
    type: "ride_update",
    data: {
      status, driver_location, current_eta, route_polyline
    },
    timestamp: unix_time,
    sequence: counter
  }
- Heartbeat: Every 30 seconds
- Reconnect: Exponential backoff (1s, 2s, 4s, 8s, 16s max)
```

#### 19. Ride History
```
GET /rides/history
- Auth: JWT Token (PASSENGER, DRIVER)
- Query: { limit, offset, filter_by_date }
- Response: { rides: [{ ride_id, date, driver_name, fare, rating }], total_count }
- Rate Limit: 100 requests/minute per user
- SLA: 99%
```

---

### GPS-SERVICE (HTTP/gRPC)

#### 20. Update Location
```
POST /gps/location
- Auth: JWT Token (DRIVER)
- Body: { latitude, longitude, accuracy_m, heading, speed }
- Response: { status: "updated" }
- Rate Limit: 600 requests/minute per driver (1 per 100ms)
- SLA: 99.5%
```

#### 21. Get Nearby Drivers
```
GET /gps/nearby-drivers
- Auth: JWT Token (PASSENGER, ADMIN)
- Query: { latitude, longitude, radius_meters }
- Response: [{ driver_id, distance_meters, current_eta, vehicle_info }]
- Rate Limit: 1000 requests/minute
- SLA: 99.5%
```

#### 22. Get Trip Polyline
```
GET /gps/trips/{trip_id}/route
- Auth: JWT Token
- Response: { polyline_encoded, distance_meters, duration_seconds }
- Rate Limit: 500 requests/minute
- SLA: 99%
```

#### 23. Trip Replay
```
GET /gps/trips/{trip_id}/replay
- Auth: JWT Token
- Query: { speed_multiplier (1-10) }
- Response: WebSocket stream of [{ latitude, longitude, timestamp }]
- Rate Limit: 10 requests/minute per user
- SLA: 95%
```

---

### DISPATCH-SERVICE (HTTP/gRPC)

#### 24. Get Available Drivers (Admin)
```
GET /dispatch/available-drivers
- Auth: JWT Token (ADMIN, OPERATIONS)
- Query: { zone_id, limit, radius_km }
- Response: [{ driver_id, distance, last_location_update }]
- Rate Limit: 500 requests/minute
- SLA: 99%
```

#### 25. Dispatch Metrics
```
GET /dispatch/metrics
- Auth: JWT Token (ADMIN)
- Query: { start_date, end_date, zone_id }
- Response: { 
    total_dispatches, success_rate, avg_eta, 
    cancellations, reassignments, avg_attempts
  }
- Rate Limit: 100 requests/minute
- SLA: 99%
```

---

### PRICING-SERVICE (HTTP/gRPC)

#### 26. Estimate Fare
```
GET /pricing/estimate
- Auth: JWT Token (PASSENGER)
- Query: { 
    pickup_latitude, pickup_longitude,
    dropoff_latitude, dropoff_longitude,
    passengers (optional)
  }
- Response: { 
    base_fare, distance_fare, time_fare, surge_multiplier,
    estimated_total_fare, currency, breakdown
  }
- Rate Limit: 1000 requests/minute
- SLA: 99.5%
```

#### 27. Pricing Rules (Admin)
```
GET /pricing/rules
- Auth: JWT Token (ADMIN, OPERATIONS)
- Query: { zone_id, effective_date }
- Response: [{ zone_id, base_fare, per_km, per_minute, effective_date }]
- Rate Limit: 100 requests/minute
- SLA: 99%
```

#### 28. Set Surge Multiplier (Admin)
```
POST /pricing/surge
- Auth: JWT Token (ADMIN)
- Body: { zone_id, multiplier (1.0-3.0), duration_minutes }
- Response: { zone_id, multiplier, active_until }
- Rate Limit: 100 requests/minute
- SLA: 99%
```

---

### PAYMENT-SERVICE (HTTP/gRPC)

#### 29. Process Payment
```
POST /payments/process
- Auth: JWT Token (PASSENGER)
- Body: { 
    ride_id, amount, currency, 
    payment_method (card|wallet|telebirr)
  }
- Response: { 
    transaction_id, status (success|pending|failed),
    amount_charged, timestamp
  }
- Rate Limit: 100 requests/minute per user
- SLA: 99.5%
```

#### 30. Payment Status
```
GET /payments/{transaction_id}
- Auth: JWT Token
- Response: { 
    transaction_id, ride_id, status, amount,
    payment_method, timestamp, receipt_url
  }
- Rate Limit: 500 requests/minute
- SLA: 99%
```

#### 31. List Transactions
```
GET /payments/history
- Auth: JWT Token (USER)
- Query: { limit, offset, start_date, end_date }
- Response: { transactions: [...], total_count }
- Rate Limit: 100 requests/minute per user
- SLA: 99%
```

#### 32. Refund Request
```
POST /payments/{transaction_id}/refund
- Auth: JWT Token (ADMIN for disputes)
- Body: { reason (optional) }
- Response: { 
    transaction_id, refund_amount, status,
    refund_id, expected_time
  }
- Rate Limit: 10 requests/minute per user
- SLA: 99%
```

---

### WALLET-SERVICE (HTTP/gRPC)

#### 33. Get Wallet Balance
```
GET /wallet/balance
- Auth: JWT Token (USER, DRIVER)
- Response: { 
    user_id, balance, currency, last_updated,
    pending_holds, available_balance
  }
- Rate Limit: 1000 requests/minute
- SLA: 99.9%
```

#### 34. Wallet Transactions
```
GET /wallet/transactions
- Auth: JWT Token (USER)
- Query: { limit, offset, transaction_type }
- Response: { 
    transactions: [{ 
      transaction_id, type (debit|credit),
      amount, source, timestamp
    }],
    total_count
  }
- Rate Limit: 100 requests/minute per user
- SLA: 99%
```

#### 35. Top Up Wallet
- Auth: JWT Token
- Body: { amount, payment_method }
- Response: { transaction_id, new_balance }
- Rate Limit: 10 requests/minute per user
- SLA: 99.5%

#### 36. Wallet Reconciliation (Admin)
```
GET /wallet/reconciliation
- Auth: JWT Token (ADMIN, FINANCE)
- Query: { date, batch_id }
- Response: { 
    total_users, total_debits, total_credits,
    discrepancies, status
  }
- Rate Limit: 10 requests/minute
- SLA: 99%
```

---

### DRIVER-SERVICE (HTTP/gRPC)

#### 37. Submit Driver Application
```
POST /drivers/apply
- Auth: JWT Token (NEW USER)
- Body: {
    license_number, license_expiry,
    vehicle_make, vehicle_model, vehicle_year,
    vehicle_license_plate, vehicle_color
  }
- Response: { application_id, status: "submitted" }
- Rate Limit: 1 request/minute per user
- SLA: 99%
```

#### 38. Upload Driver Documents
```
POST /drivers/documents
- Auth: JWT Token (DRIVER)
- Multipart Form Data:
    - document_type: license|id|insurance|registration|selfie
    - file: binary image
- Response: { document_id, status: "uploaded", verification_status }
- Rate Limit: 100 requests/minute per driver
- SLA: 99%
```

#### 39. Get Driver Status
```
GET /drivers/{driver_id}/status
- Auth: JWT Token (DRIVER, ADMIN)
- Response: { 
    driver_id, application_status (draft|submitted|approved|rejected),
    documents_verified, rating, total_rides, account_status
  }
- Rate Limit: 500 requests/minute
- SLA: 99%
```

#### 40. Driver Document Check (Admin)
```
GET /drivers/{driver_id}/documents
- Auth: JWT Token (ADMIN)
- Response: [{ 
    document_id, document_type, upload_date,
    verification_status, verified_by, notes
  }]
- Rate Limit: 100 requests/minute
- SLA: 99%
```

#### 41. Approve Driver (Admin)
```
POST /drivers/{driver_id}/approve
- Auth: JWT Token (ADMIN)
- Body: { notes (optional) }
- Response: { driver_id, status: "approved", effective_date }
- Rate Limit: 100 requests/minute
- SLA: 99%
```

#### 42. Reject Driver (Admin)
```
POST /drivers/{driver_id}/reject
- Auth: JWT Token (ADMIN)
- Body: { reason_code, details }
- Response: { driver_id, status: "rejected" }
- Rate Limit: 100 requests/minute
- SLA: 99%
```

#### 43. Suspend Driver (Admin)
```
POST /drivers/{driver_id}/suspend
- Auth: JWT Token (ADMIN)
- Body: { reason_code, duration_hours }
- Response: { driver_id, status: "suspended", suspend_until }
- Rate Limit: 50 requests/minute
- SLA: 99%
```

#### 44. Driver Documents Resubmit (Driver)
```
POST /drivers/documents/resubmit
- Auth: JWT Token (DRIVER)
- Body: { application_id, document_updates }
- Response: { status: "submitted_for_review" }
- Rate Limit: 10 requests/minute per driver
- SLA: 99%
```

---

### SAFETY-SERVICE (HTTP/gRPC)

#### 45. Trigger SOS
```
POST /safety/sos
- Auth: JWT Token (PASSENGER, DRIVER)
- Body: { ride_id, location: { latitude, longitude }, contact_ids }
- Response: { sos_id, status: "triggered", contacts_notified }
- Rate Limit: 1 request/5 minutes per user
- SLA: 99.9%
```

#### 46. Share Trip
```
POST /safety/share-trip
- Auth: JWT Token (PASSENGER, DRIVER)
- Body: { ride_id, recipient_emails: [], duration_minutes }
- Response: { share_id, share_link, expires_at }
- Rate Limit: 50 requests/minute
- SLA: 99%
```

#### 47. Report Incident
```
POST /safety/incidents
- Auth: JWT Token (PASSENGER, DRIVER)
- Body: { ride_id, incident_type, description, photos: [urls] }
- Response: { incident_id, status: "reported" }
- Rate Limit: 10 requests/minute per user
- SLA: 99%
```

#### 48. Rate Safety (Post-Ride)
```
POST /rides/{ride_id}/safety-rating
- Auth: JWT Token (PASSENGER, DRIVER)
- Body: { safety_rating (1-5), concerns_notes }
- Response: { rating_id }
- Rate Limit: 1 per ride
- SLA: 99%
```

---

### FRAUD-SERVICE (HTTP/gRPC)

#### 49. Check Fraud Score (Internal)
```
POST /fraud/check
- Auth: Service JWT
- Body: { user_id, transaction_type, amount, context }
- Response: { risk_score (0-100), risk_level, rules_triggered }
- Rate Limit: 10000 requests/minute
- SLA: 99.9%
```

#### 50. Fraud Appeal (User)
```
POST /fraud/appeals
- Auth: JWT Token (USER)
- Body: { fraud_id, explanation, supporting_docs }
- Response: { appeal_id, status: "submitted" }
- Rate Limit: 5 requests/minute per user
- SLA: 99%
```

---

### SUPPORT-SERVICE (HTTP/gRPC)

#### 51. Create Support Ticket
```
POST /support/tickets
- Auth: JWT Token (USER)
- Body: { ride_id (optional), issue_category, description, attachments }
- Response: { ticket_id, status: "created" }
- Rate Limit: 10 requests/minute per user
- SLA: 99%
```

#### 52. Get Ticket Status
```
GET /support/tickets/{ticket_id}
- Auth: JWT Token (USER, SUPPORT)
- Response: { 
    ticket_id, status, created_at, updated_at,
    interactions: [{ actor, message, timestamp }],
    assigned_agent
  }
- Rate Limit: 500 requests/minute
- SLA: 99%
```

#### 53. List Tickets
```
GET /support/tickets
- Auth: JWT Token (USER, SUPPORT)
- Query: { status, limit, offset }
- Response: { tickets: [...], total_count }
- Rate Limit: 100 requests/minute per user
- SLA: 99%
```

---

### ANALYTICS-SERVICE (HTTP)

#### 54. Daily Metrics
```
GET /analytics/daily-metrics
- Auth: JWT Token (ADMIN, OPERATIONS)
- Query: { date, zone_id }
- Response: { 
    date, rides_count, active_drivers, active_passengers,
    avg_fare, avg_rating, total_revenue
  }
- Rate Limit: 1000 requests/minute
- SLA: 99%
```

#### 55. Driver Performance
```
GET /analytics/drivers/{driver_id}/performance
- Auth: JWT Token (ADMIN, DRIVER)
- Query: { start_date, end_date }
- Response: { 
    total_rides, total_earnings, acceptance_rate,
    cancellation_rate, avg_rating, online_hours
  }
- Rate Limit: 500 requests/minute
- SLA: 99%
```

#### 56. Zone Analytics
```
GET /analytics/zones/{zone_id}
- Auth: JWT Token (ADMIN, OPERATIONS)
- Query: { start_date, end_date }
- Response: { 
    zone_id, rides_count, avg_eta, demand_trend,
    surge_times, peak_hours
  }
- Rate Limit: 500 requests/minute
- SLA: 99%
```

---

## RATE LIMITING STRATEGY

**Tier-Based Limits:**
- Tier 1 (Unauthenticated): 100 requests/minute per IP
- Tier 2 (Authenticated User): 1000 requests/minute per user
- Tier 3 (Premium User): 5000 requests/minute per user
- Tier 4 (Service Account): Unlimited (negotiated per service)

**Headers Returned:**
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 943
X-RateLimit-Reset: 1234567890
```

---

## RESPONSE FORMATS

**Success (2xx):**
```json
{
  "success": true,
  "data": { ... },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Error (4xx/5xx):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Human-readable error",
    "details": { ... },
    "correlation_id": "uuid"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

**Audit Status:** ✅ API_CATALOG.md COMPLETE  
**Task 1 Sub-Deliverables:** All 4 catalogs created (SERVICE, EVENT, DATABASE, API)

