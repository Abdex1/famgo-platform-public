# 🔧 SERVICE BOOTSTRAP PATTERN
## Extracted from uber-master, Adapted for FamGo

**Status:** Pattern 2/8 - EXTRACTED AND DOCUMENTED  
**Source:** uber-master (`services/user-service/cmd/main.go`)  
**Adoption Level:** Category A (Directly Adopt)  
**Risk Level:** LOW

---

## PATTERN OVERVIEW

### Service Initialization Sequence

```
1. Load Configuration
   └── From environment variables

2. Setup Logging
   └── Structured logging initialization

3. Connect to Database
   └── PostgreSQL connection pool

4. Run Migrations
   └── Database schema updates

5. Connect to Cache
   └── Redis client setup

6. Setup Kafka (if needed)
   └── Producer/consumer initialization

7. Initialize Services
   └── Create service layer

8. Setup HTTP Router
   └── Register handlers and middleware

9. Start Health Checks
   └── /health and /ready endpoints

10. Start HTTP Server
    └── Listen and serve

11. Graceful Shutdown
    └── Signal handling and cleanup
```

---

## IMPLEMENTATION TEMPLATE

```go
// cmd/main.go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/famgo/shared/pkg/config"
    "github.com/famgo/shared/pkg/db"
    "github.com/famgo/shared/pkg/redis"
    "github.com/famgo/shared/pkg/kafka"
    "github.com/famgo/shared/pkg/http"
    "github.com/famgo/shared/pkg/observability"
    
    "{service}/internal/handler"
    "{service}/internal/repository"
    "{service}/internal/service"
    "{service}/config"
)

func main() {
    ctx := context.Background()
    
    // STEP 1: Load configuration
    cfg := config.LoadConfig()
    
    // STEP 2: Setup logging
    logger := observability.NewLogger(cfg.LogLevel)
    logger.Info("Starting service", "service", cfg.ServiceName)
    
    // STEP 3: Connect to database
    dbPool, err := db.NewPool(ctx, cfg.DatabaseURL)
    if err != nil {
        logger.Error("Failed to connect to database", "error", err)
        os.Exit(1)
    }
    defer dbPool.Close()
    logger.Info("Connected to database")
    
    // STEP 4: Run migrations
    if err := db.RunMigrations(ctx, dbPool, "migrations"); err != nil {
        logger.Error("Failed to run migrations", "error", err)
        os.Exit(1)
    }
    logger.Info("Migrations completed")
    
    // STEP 5: Connect to Redis
    redisClient, err := redis.NewClient(ctx, cfg.RedisAddr)
    if err != nil {
        logger.Error("Failed to connect to Redis", "error", err)
        os.Exit(1)
    }
    defer redisClient.Close()
    logger.Info("Connected to Redis")
    
    // STEP 6: Setup Kafka (if needed)
    var kafkaProducer kafka.Producer
    if cfg.KafkaEnabled {
        kafkaProducer, err = kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaUsername, cfg.KafkaPassword)
        if err != nil {
            logger.Error("Failed to setup Kafka producer", "error", err)
            os.Exit(1)
        }
        defer kafkaProducer.Close()
        logger.Info("Kafka producer initialized")
    }
    
    // STEP 7: Initialize service layer
    repo := repository.New(dbPool)
    svc := service.New(repo, kafkaProducer, redisClient, logger)
    
    // STEP 8: Setup HTTP router
    router := http.NewRouter()
    
    // Register health checks
    router.Get("/health", http.HealthCheck)
    router.Get("/ready", http.ReadinessCheck(dbPool, redisClient))
    
    // Register service handlers
    handler.RegisterRoutes(router, svc)
    
    // STEP 9: Metrics endpoint
    router.Get("/metrics", observability.MetricsHandler)
    
    // STEP 10: Start HTTP server
    server := &http.Server{
        Addr:    fmt.Sprintf(":%s", cfg.ServicePort),
        Handler: router,
    }
    
    // Start server in goroutine
    go func() {
        logger.Info("Server starting", "port", cfg.ServicePort)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("Server error", "error", err)
        }
    }()
    
    // STEP 11: Graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigChan
    logger.Info("Shutting down gracefully")
    
    // Shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(shutdownCtx); err != nil {
        logger.Error("Server shutdown error", "error", err)
    }
    
    logger.Info("Server stopped")
}
```

---

## KEY COMPONENTS

### Health Check
```go
func HealthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
```

### Readiness Check
```go
func ReadinessCheck(db *pgxpool.Pool, cache *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check critical dependencies
        if err := db.Ping(r.Context()); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }
        if err := cache.Ping(r.Context()).Err(); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
    }
}
```

### Graceful Shutdown
```go
// Listen for signals
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// Wait for signal
<-sigChan

// Shutdown with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
server.Shutdown(ctx)
```

---

**Pattern 2 Status:** READY FOR USE  
**All Services:** Must follow this bootstrap sequence  
**Risk:** LOW

---
