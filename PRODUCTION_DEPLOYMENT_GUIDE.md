# 📚 PRODUCTION DEPLOYMENT & OPERATIONS GUIDE

**Platform:** FamGo Microservices  
**Version:** 1.0.0  
**Status:** Production Ready  
**Last Updated:** Weeks 3-4 Completion

---

## TABLE OF CONTENTS

1. Deployment Guide
2. Observability Setup
3. Security Configuration
4. Troubleshooting Guide
5. Incident Response
6. Scaling Guidelines
7. Maintenance Procedures

---

## 1. DEPLOYMENT GUIDE

### Prerequisites

- Kubernetes 1.24+
- Helm 3.0+
- Docker 20.10+
- PostgreSQL 14+
- Redis 7.0+
- Prometheus, Jaeger, Loki (monitoring stack)

### Services Deployment Order

```
1. PostgreSQL (persistence layer)
2. Redis (caching layer)
3. Kafka (event bus)
4. Auth Service (security foundation)
5. User Service (user management)
6. GPS Service (location tracking)
7. Pricing Service (fare calculation)
8. Dispatch Service (driver assignment)
9. Ride Service (main orchestration)
10. Payment Service (financial transactions)
```

### Deployment via Kubernetes

```bash
# Create namespace
kubectl create namespace famgo

# Deploy databases
kubectl apply -f deployments/databases/postgres.yaml
kubectl apply -f deployments/databases/redis.yaml

# Deploy services (follow order above)
kubectl apply -f services/auth-service/deployments/kubernetes.yaml
kubectl apply -f services/user-service/deployments/kubernetes.yaml
kubectl apply -f services/gps-service/deployments/kubernetes.yaml
kubectl apply -f services/ride-service/deployments/kubernetes.yaml

# Verify rollout
kubectl rollout status deployment/ride-service -n famgo --timeout=5m

# Port forward for testing
kubectl port-forward svc/ride-service 8080:80 -n famgo
```

### Health Check Verification

```bash
# Check liveness (service is running)
curl http://localhost:8080/health

# Check readiness (ready to serve traffic)
curl http://localhost:8080/ready

# Expected responses: 200 OK with {"status": "healthy"}
```

---

## 2. OBSERVABILITY SETUP

### Prometheus Configuration

**Scrape interval:** 15 seconds  
**Retention:** 15 days  
**Metrics endpoint:** `/metrics` on each service

```yaml
# prometheus.yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'ride-service'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

### Jaeger Trace Configuration

**Sampling:** 10% of requests  
**Collector endpoint:** `jaeger-collector:14250`  
**Query UI:** `http://jaeger-ui:16686`

### Loki Log Aggregation

**Log format:** JSON structured logs  
**Retention:** 30 days  
**Query language:** LogQL

Sample structured log:
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "service": "ride-service",
  "operation": "create_ride",
  "user_id": "user123",
  "trace_id": "abc123def456",
  "duration_ms": 45,
  "message": "Ride created successfully"
}
```

### Grafana Dashboards

**Pre-built dashboards:**
1. Request Performance (latency p50, p95, p99)
2. Ride Metrics (created, completed, cancelled, active)
3. gRPC Calls (rate, latency, errors by service)
4. Circuit Breaker Status (health of cross-service calls)
5. Resource Usage (CPU, memory, disk, pod count)

**Access:** `http://grafana:3000`  
**Default credentials:** admin/admin

---

## 3. SECURITY CONFIGURATION

### JWT Token Validation

**Token endpoint:** `auth-service:50051` (gRPC)  
**Token format:** Bearer <JWT>  
**Validation:** Signature, expiration, claims

```go
// Example validation
token := extractBearerToken(request)
authCtx, err := authMiddleware.ValidateToken(ctx, token)
if err != nil {
  return 401, "Unauthorized"
}
```

### RBAC Roles

**Defined Roles:**
- `PASSENGER`: Create rides, view own rides
- `DRIVER`: Accept rides, update location, start/complete rides
- `DISPATCHER`: Assign drivers, cancel assignments
- `ADMIN`: Full access to all resources

**Authorization checks:**
- All endpoints require valid JWT token
- Endpoint access controlled by required roles
- Audit logging on authorization failures

### Input Validation Rules

**Ride Creation:**
- Passenger ID: Non-empty, <100 chars
- Coordinates: Valid latitude (-90 to 90), longitude (-180 to 180)
- Pickup ≠ Dropoff location

**Fare Amount:**
- Minimum: $1.00
- Maximum: $10,000.00

**IDs:**
- Format: Non-empty, <100 chars
- Pattern: Alphanumeric with hyphens allowed

### Audit Logging

**Logged Events:**
- User authentication (success/failure)
- Authorization decisions (allow/deny)
- Sensitive operations (create, update, delete)
- Security events (invalid input, rate limits)

**Audit log location:** `/var/log/audit.log` (JSON format)

---

## 4. TROUBLESHOOTING GUIDE

### Common Issues

#### Service Won't Start

**Symptoms:** Pod stuck in CrashLoopBackOff

**Debug Steps:**
```bash
# Check pod logs
kubectl logs <pod-name> -n famgo

# Check pod description
kubectl describe pod <pod-name> -n famgo

# Check events
kubectl get events -n famgo --sort-by='.lastTimestamp'
```

**Common causes:**
- Database connection failed → check DB pod is running
- Redis connection failed → check Redis pod is running
- Invalid environment variables → check ConfigMap

#### High Latency

**Symptoms:** Request latency > 1 second

**Debug Steps:**
```bash
# Check Prometheus metrics
# Query: histogram_quantile(0.95, http_request_duration_seconds_bucket)

# Check circuit breaker status
# Query: circuit_breaker_status

# Check service dependencies
# Query: rate(grpc_call_count[1m]) by (service)
```

**Common causes:**
- Circuit breaker open → downstream service failing
- Database slow → check DB indexes, query performance
- Network latency → check K8s network policies

#### High Error Rate

**Symptoms:** Error rate > 5%

**Debug Steps:**
```bash
# Check error types
kubectl logs <pod-name> -n famgo | grep ERROR

# Check Prometheus
# Query: rate(http_request_errors_total[1m])

# Check Loki logs
# Query: {service="ride-service"} | json | level="ERROR"
```

**Common causes:**
- Invalid input → check input validation rules
- Database constraint violation → check data consistency
- Service dependency down → check cross-service calls

#### Memory Leak

**Symptoms:** Memory usage increasing over time

**Debug Steps:**
```bash
# Check memory usage trend
# Query: container_memory_usage_bytes{pod=~'ride-service.*'}

# Check goroutine count
# Query: go_goroutines{job="ride-service"}

# Profile memory
curl http://localhost:8080/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

**Common causes:**
- Unclosed connections → check connection pooling
- Memory leak in library → update dependencies
- Large cache → check cache eviction policy

---

## 5. INCIDENT RESPONSE

### Critical Incident (Service Down)

**1. Immediate Actions (< 5 min)**
```bash
# Verify pod is running
kubectl get pods -n famgo -l app=ride-service

# Check service endpoint
kubectl get endpoints ride-service -n famgo

# Check recent logs
kubectl logs ride-service-xxx -n famgo --tail=100
```

**2. Root Cause Analysis (< 15 min)**
```bash
# Database health
kubectl exec -it postgres-pod -- psql -U postgres -c "SELECT 1"

# Redis health
kubectl exec -it redis-pod -- redis-cli ping

# Network connectivity
kubectl run -it --rm debug --image=nicolaka/netshoot -- /bin/bash
# Inside: nslookup ride-service
# Inside: curl http://ride-service:8080/ready
```

**3. Recovery (< 30 min)**
```bash
# Option 1: Restart pod
kubectl rollout restart deployment/ride-service -n famgo

# Option 2: Rollback to previous version
kubectl rollout undo deployment/ride-service -n famgo

# Option 3: Scale down and up
kubectl scale deployment ride-service --replicas=0 -n famgo
kubectl scale deployment ride-service --replicas=3 -n famgo
```

### High Latency Incident

**1. Immediate Actions**
```bash
# Check circuit breaker status
# Query in Prometheus: circuit_breaker_status{service="pricing"}

# Check gRPC call latency
# Query: histogram_quantile(0.99, grpc_call_duration_seconds_bucket)

# Identify slow service
# Query: rate(grpc_call_count[1m]) by (service, status)
```

**2. Mitigation**
```bash
# Increase timeouts temporarily
kubectl set env deployment/ride-service \
  PRICING_TIMEOUT_MS=10000 -n famgo

# Scale down to force load balancing
kubectl scale deployment/ride-service --replicas=1 -n famgo
# Wait 2 min for traffic to balance

# Scale back up
kubectl scale deployment/ride-service --replicas=3 -n famgo
```

### Data Integrity Issue

**1. Verification**
```bash
# Check ride count consistency
SELECT COUNT(*) FROM rides;
SELECT COUNT(*) FROM ride_status_history;

# Check for orphaned records
SELECT * FROM rides WHERE driver_id NOT IN (SELECT id FROM drivers);
```

**2. Resolution**
```bash
# Backup data
pg_dump ride_db > backup_$(date +%Y%m%d).sql

# Run integrity check
psql ride_db -c "SELECT pg_catalog.pg_constraint_check('rides')"

# Repair if needed
# (Contact data team for guidance)
```

---

## 6. SCALING GUIDELINES

### Horizontal Scaling

**When to scale:**
- CPU usage > 70% consistently
- Memory usage > 80% consistently
- Request latency > 500ms p95
- Request queue depth > 100

**Scale up:**
```bash
kubectl scale deployment/ride-service --replicas=10 -n famgo
```

**Scale down:**
```bash
kubectl scale deployment/ride-service --replicas=3 -n famgo
```

### Vertical Scaling

**Current resource limits:**
- CPU: 1000m (1 CPU)
- Memory: 512Mi

**Increase:**
```bash
kubectl set resources deployment/ride-service \
  --limits=cpu=2000m,memory=1Gi -n famgo
```

### Auto-Scaling

**HPA Configuration:**
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ride-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ride-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

---

## 7. MAINTENANCE PROCEDURES

### Database Maintenance

**Daily:**
```bash
# Check disk space
df -h

# Monitor slow queries
SELECT query, mean_exec_time FROM pg_stat_statements 
ORDER BY mean_exec_time DESC LIMIT 10;
```

**Weekly:**
```bash
# Vacuum and analyze
VACUUM ANALYZE;

# Update statistics
ANALYZE;
```

**Monthly:**
```bash
# Reindex
REINDEX DATABASE ride_db;

# Check table size
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename))
FROM pg_tables 
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### Backup & Recovery

**Automated backups:**
- Frequency: Daily at 2 AM UTC
- Retention: 30 days
- Location: S3 bucket

**Manual backup:**
```bash
pg_dump ride_db | gzip > backup_$(date +%Y%m%d).sql.gz
aws s3 cp backup_*.sql.gz s3://backups/
```

**Recovery:**
```bash
gunzip < backup_20240115.sql.gz | psql ride_db
```

### Dependency Updates

**Check for updates:**
```bash
go list -u -m all

# Security vulnerabilities
go list -m -json all | nancy sleuth
```

**Update safely:**
```bash
# Test in staging first
go get -u github.com/package/name@v1.2.3

# Run tests
go test ./...

# Deploy to staging
kubectl apply -f deployments/staging/

# Monitor for issues (24 hours)
# Then deploy to production
kubectl apply -f deployments/kubernetes.yaml
```

### Rolling Updates

```bash
# Gradually replace pods
kubectl set image deployment/ride-service \
  ride-service=ride-service:v1.2.0 \
  --record -n famgo

# Monitor rollout
kubectl rollout status deployment/ride-service -n famgo

# If issues, rollback
kubectl rollout undo deployment/ride-service -n famgo
```

---

## SUPPORT & CONTACTS

**On-Call:** See Slack #oncall-schedule  
**Escalation:** #incident-commander  
**Documentation:** https://wiki.famgo.internal/  
**Runbooks:** https://runbooks.famgo.internal/

---

**Last Updated:** Days 1-10 Complete  
**Next Review:** Monthly
