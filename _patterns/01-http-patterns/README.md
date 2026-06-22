# 🔌 HTTP HANDLER PATTERNS
## Extracted from uber-master, Adapted for FamGo

**Status:** Pattern 1/8 - EXTRACTED AND DOCUMENTED  
**Source:** uber-master (`services/user-service/internal/handler/`)  
**Adoption Level:** Category A (Directly Adopt)  
**Risk Level:** LOW  
**Applicable To:** All 19 FamGo services

---

## PATTERN OVERVIEW

### What This Pattern Provides

```
Chi Router Setup
├── HTTP handler registration
├── Middleware composition
├── Request routing
├── Route parameter extraction
└── Handler error handling

Middleware Stack
├── Request validation
├── Authentication
├── Error handling
├── CORS support
├── Request logging
└── Timeout enforcement

Response Formatting
├── Success response envelope
├── Error response envelope
├── Status code mapping
└── JSON serialization
```

### Why Extract This

```
✅ Proven in production (uber-master uses it)
✅ Clean code organization
✅ Reusable middleware composition
✅ Consistent error handling
✅ Standard response format
✅ Easy to extend per service
```

---

## PATTERN SPECIFICATION

### 1. Chi Router Setup

**From Uber:**
```go
// From uber-master: services/user-service/internal/handler/

import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
    "github.com/go-chi/middleware"
)

// Create router with standard middleware
func NewRouter() *chi.Mux {
    r := chi.NewRouter()
    
    // Global middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(30 * time.Second))
    
    // CORS
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Content-Length"},
        MaxAge:           300,
        AllowCredentials: true,
    }))
    
    return r
}
```

**Adapt for FamGo:**
```go
// For FamGo: shared/pkg/http/router.go

package http

import (
    "time"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
    "github.com/go-chi/middleware"
)

// NewRouter creates a standard FamGo HTTP router
func NewRouter() *chi.Mux {
    r := chi.NewRouter()
    
    // Standard middleware stack
    r.Use(middleware.RequestID)           // Add X-Request-ID
    r.Use(middleware.RealIP)              // Get real IP
    r.Use(middleware.Logger)              // Log requests
    r.Use(middleware.Recoverer)           // Recover from panics
    r.Use(middleware.Timeout(30 * time.Second))
    
    // CORS configuration
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},  // Configure per environment
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept"},
        ExposedHeaders:   []string{"Content-Length", "X-Request-ID", "X-Trace-ID"},
        MaxAge:           300,
        AllowCredentials: true,
    }))
    
    return r
}

// Health check endpoint (for all services)
func HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// Readiness check endpoint (for all services)
func ReadinessCheck(db *pgxpool.Pool, cache *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check database
        if err := db.Ping(r.Context()); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            json.NewEncoder(w).Encode(map[string]string{"error": "database unavailable"})
            return
        }
        
        // Check Redis
        if err := cache.Ping(r.Context()).Err(); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            json.NewEncoder(w).Encode(map[string]string{"error": "cache unavailable"})
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
    }
}
```

---

### 2. Middleware Stack

#### 2.1 Authentication Middleware

**Pattern:**
```go
// shared/pkg/http/middleware_auth.go

package http

import (
    "context"
    "net/http"
    "strings"
    "github.com/famgo/shared/pkg/security"
)

// AuthMiddleware enforces authentication
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            Error(w, "", http.StatusUnauthorized, "UNAUTHENTICATED", "Missing Authorization header", "")
            return
        }
        
        // Parse Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            Error(w, "", http.StatusUnauthorized, "INVALID_TOKEN", "Invalid Authorization format", "")
            return
        }
        
        token := parts[1]
        
        // Verify token
        claims, err := security.VerifyToken(r.Context(), token)
        if err != nil {
            Error(w, "", http.StatusUnauthorized, "INVALID_TOKEN", "Token verification failed", err.Error())
            return
        }
        
        // Store claims in context
        ctx := context.WithValue(r.Context(), "claims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// OptionalAuthMiddleware allows requests with or without auth
func OptionalAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        
        if authHeader != "" {
            parts := strings.Split(authHeader, " ")
            if len(parts) == 2 && parts[0] == "Bearer" {
                claims, err := security.VerifyToken(r.Context(), parts[1])
                if err == nil {
                    ctx := context.WithValue(r.Context(), "claims", claims)
                    next.ServeHTTP(w, r.WithContext(ctx))
                    return
                }
            }
        }
        
        // Continue without auth if header invalid or missing
        next.ServeHTTP(w, r)
    })
}
```

#### 2.2 Validation Middleware

**Pattern:**
```go
// shared/pkg/http/middleware_validation.go

package http

import (
    "encoding/json"
    "net/http"
)

// ValidateJSON middleware validates request body is valid JSON
func ValidateJSON(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.ContentLength == 0 {
            next.ServeHTTP(w, r)
            return
        }
        
        // Check Content-Type
        if ct := r.Header.Get("Content-Type"); ct != "application/json" {
            Error(w, "", http.StatusBadRequest, "INVALID_CONTENT_TYPE", "Content-Type must be application/json", "")
            return
        }
        
        // Validate JSON structure (without parsing)
        decoder := json.NewDecoder(r.Body)
        var body map[string]interface{}
        if err := decoder.Decode(&body); err != nil {
            Error(w, "", http.StatusBadRequest, "INVALID_JSON", "Request body must be valid JSON", err.Error())
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

#### 2.3 Error Handling Middleware

**Pattern:**
```go
// shared/pkg/http/middleware_error.go

package http

import (
    "net/http"
    "encoding/json"
)

type ErrorResponse struct {
    Success bool        `json:"success"`
    TraceID string      `json:"trace_id"`
    Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// Error sends standardized error response
func Error(w http.ResponseWriter, traceID, statusCode, code, message, details string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := ErrorResponse{
        Success: false,
        TraceID: traceID,
        Error: &ErrorDetail{
            Code:    code,
            Message: message,
            Details: details,
        },
    }
    
    json.NewEncoder(w).Encode(response)
}

// Success sends standardized success response
func Success(w http.ResponseWriter, traceID string, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    response := map[string]interface{}{
        "success":  true,
        "trace_id": traceID,
        "data":     data,
    }
    
    json.NewEncoder(w).Encode(response)
}
```

---

### 3. Handler Pattern

**From Uber:**
```go
// From uber-master: services/user-service/internal/handler/user.go

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    
    user, err := h.service.GetUser(r.Context(), userID)
    if err != nil {
        // Handle error
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}
```

**Adapt for FamGo (Using Response Envelope):**
```go
// For FamGo: services/{service}/internal/handler/{domain}.go

package handler

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/famgo/shared/pkg/http"
)

type Handler struct {
    service *service.Service
    logger  log.Logger
}

// GetEntity retrieves an entity by ID
func (h *Handler) GetEntity(w http.ResponseWriter, r *http.Request) {
    // Get trace ID from context (added by middleware)
    traceID := r.Header.Get("X-Request-ID")
    
    // Extract ID from URL parameter
    entityID := chi.URLParam(r, "id")
    if entityID == "" {
        http.Error(w, traceID, http.StatusBadRequest, "INVALID_INPUT", "Entity ID required", "")
        return
    }
    
    // Get from service
    entity, err := h.service.GetEntity(r.Context(), entityID)
    if err != nil {
        // Log error
        h.logger.Error("Failed to get entity", "error", err, "entity_id", entityID)
        
        // Return error response
        if errors.Is(err, service.ErrNotFound) {
            http.Error(w, traceID, http.StatusNotFound, "NOT_FOUND", "Entity not found", "")
        } else {
            http.Error(w, traceID, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve entity", err.Error())
        }
        return
    }
    
    // Return success response
    http.Success(w, traceID, http.StatusOK, entity)
}

// CreateEntity creates a new entity
func (h *Handler) CreateEntity(w http.ResponseWriter, r *http.Request) {
    traceID := r.Header.Get("X-Request-ID")
    
    // Parse request body
    var req CreateEntityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, traceID, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", err.Error())
        return
    }
    
    // Validate input
    if err := req.Validate(); err != nil {
        http.Error(w, traceID, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), "")
        return
    }
    
    // Get claims from context (added by auth middleware)
    claims, ok := r.Context().Value("claims").(*security.Claims)
    if !ok {
        http.Error(w, traceID, http.StatusUnauthorized, "UNAUTHENTICATED", "Not authenticated", "")
        return
    }
    
    // Create entity
    entity, err := h.service.CreateEntity(r.Context(), req, claims.UserID)
    if err != nil {
        h.logger.Error("Failed to create entity", "error", err, "user_id", claims.UserID)
        http.Error(w, traceID, http.StatusInternalServerError, "CREATION_FAILED", "Failed to create entity", err.Error())
        return
    }
    
    // Return created response
    http.Success(w, traceID, http.StatusCreated, entity)
}
```

---

## USAGE GUIDE

### For Each FamGo Service

**Step 1: Import Standard Router**
```go
import "github.com/famgo/shared/pkg/http"

// In main.go or service setup
router := http.NewRouter()
```

**Step 2: Register Health Checks**
```go
router.Get("/health", http.HealthCheck)
router.Get("/ready", http.ReadinessCheck(db, redis))
```

**Step 3: Register Routes with Middleware**
```go
// Public endpoints (no auth required)
router.Post("/register", handler.Register)
router.Post("/login", handler.Login)

// Protected endpoints (auth required)
router.Group(func(r chi.Router) {
    r.Use(http.AuthMiddleware)
    
    r.Get("/profile", handler.GetProfile)
    r.Put("/profile", handler.UpdateProfile)
    r.Delete("/profile", handler.DeleteProfile)
})

// Admin endpoints (admin auth required)
router.Group(func(r chi.Router) {
    r.Use(http.AuthMiddleware)
    r.Use(AdminAuthMiddleware)
    
    r.Get("/admin/stats", handler.GetStats)
})
```

**Step 4: Handler Implementation**
```go
// Follow pattern above - always use http.Success() and http.Error()
// Always extract traceID from context
// Always validate input
// Always use proper status codes
```

---

## EXAMPLES

### Example 1: Simple GET Endpoint

```go
func (h *Handler) GetDriver(w http.ResponseWriter, r *http.Request) {
    traceID := r.Header.Get("X-Request-ID")
    driverID := chi.URLParam(r, "id")
    
    if driverID == "" {
        http.Error(w, traceID, http.StatusBadRequest, "INVALID_INPUT", "Driver ID required", "")
        return
    }
    
    driver, err := h.service.GetDriver(r.Context(), driverID)
    if err != nil {
        if errors.Is(err, service.ErrNotFound) {
            http.Error(w, traceID, http.StatusNotFound, "NOT_FOUND", "Driver not found", "")
        } else {
            http.Error(w, traceID, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get driver", err.Error())
        }
        return
    }
    
    http.Success(w, traceID, http.StatusOK, driver)
}
```

### Example 2: Protected POST Endpoint

```go
func (h *Handler) CreateRide(w http.ResponseWriter, r *http.Request) {
    traceID := r.Header.Get("X-Request-ID")
    
    // Parse request
    var req CreateRideRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, traceID, http.StatusBadRequest, "INVALID_JSON", "Invalid request", err.Error())
        return
    }
    
    // Validate
    if err := req.Validate(); err != nil {
        http.Error(w, traceID, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), "")
        return
    }
    
    // Get authenticated user
    claims := r.Context().Value("claims").(*security.Claims)
    
    // Create ride
    ride, err := h.service.CreateRide(r.Context(), req, claims.UserID)
    if err != nil {
        http.Error(w, traceID, http.StatusInternalServerError, "CREATION_FAILED", "Failed to create ride", err.Error())
        return
    }
    
    http.Success(w, traceID, http.StatusCreated, ride)
}
```

---

## WHEN TO USE THIS PATTERN

✅ **Use this for:**
- All HTTP endpoints
- All services needing REST API
- Authentication/authorization
- Error handling
- Response formatting

❌ **Don't use this for:**
- WebSocket connections (different pattern)
- gRPC services (different protocol)
- Internal service-to-service calls (use events)

---

## INTEGRATION POINTS

**Must integrate with:**
```
✅ shared/pkg/security/         - Authentication
✅ shared/pkg/observability/    - Logging & metrics
✅ shared/pkg/errors/           - Error types
```

**May extend with:**
```
✅ Rate limiting middleware
✅ Caching middleware
✅ Circuit breaker middleware
✅ Service-specific middleware
```

---

## TESTING

**Unit test example:**
```go
func TestGetDriver(t *testing.T) {
    // Setup
    mockService := &MockService{}
    handler := &Handler{service: mockService}
    
    // Request
    req := httptest.NewRequest("GET", "/drivers/123", nil)
    w := httptest.NewRecorder()
    
    // Execute
    handler.GetDriver(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    
    var resp map[string]interface{}
    json.NewDecoder(w.Body).Decode(&resp)
    assert.True(t, resp["success"].(bool))
}
```

---

**Status:** Pattern 1 COMPLETE  
**Ready for:** All services in Week 1+  
**Dependencies:** shared/pkg/http module  
**Risk Level:** LOW (proven pattern)

---
