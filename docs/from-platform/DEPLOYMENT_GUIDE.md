# FamGo Platform - Deployment Guide

## Prerequisites

- Docker & Docker Compose
- Kubernetes cluster (1.24+)
- Terraform (1.5+)
- kubectl
- AWS CLI (for AWS deployment)

## Local Development Setup

### 1. Clone Repository
```bash
git clone https://github.com/famgo/platform.git
cd FamGo-platform
```

### 2. Start Services
```bash
docker-compose up -d

# Verify services
docker-compose ps
```

### 3. Initialize Database
```bash
psql -h localhost -U famgo -d famgo_platform -f database/migrations/001_initial_schema.sql
# ... run all migrations
```

### 4. Start Backend Services
```bash
cd backend
go mod download
go run cmd/api-gateway/main.go
```

### 5. Build Mobile Apps
```bash
# Rider App
cd mobile/flutter-passenger-app
flutter pub get
flutter run -d <device-id>

# Driver App
cd mobile/flutter-driver-app
flutter pub get
flutter run -d <device-id>
```

### 6. Run Admin Dashboard
```bash
cd web/admin-dashboard
npm install
npm run dev
```

## Docker Deployment

### Build Images
```bash
# Build all services
docker build -f backend/docker/Dockerfile.gateway -t famgo/gateway:latest .
docker build -f backend/docker/Dockerfile.ride -t famgo/ride-service:latest .
# ... for each service
```

### Deploy with Docker Compose
```bash
docker-compose -f docker-compose.yml up -d
```

## Kubernetes Deployment

### 1. Create Namespace
```bash
kubectl apply -f k8s/namespace.yaml
```

### 2. Create Secrets
```bash
kubectl create secret generic famgo-secrets \
  --from-literal=database-url=postgresql://... \
  --from-literal=jwt-secret=your-secret \
  -n famgo
```

### 3. Deploy Services
```bash
kubectl apply -f k8s/ride-service.yaml
kubectl apply -f k8s/driver-service.yaml
# ... for each service
```

### 4. Verify Deployment
```bash
kubectl get pods -n famgo
kubectl get svc -n famgo
kubectl logs -f deployment/ride-service -n famgo
```

## AWS Deployment with Terraform

### 1. Configure AWS
```bash
aws configure
export AWS_REGION=us-east-1
```

### 2. Initialize Terraform
```bash
cd infrastructure/terraform
terraform init
```

### 3. Plan Infrastructure
```bash
terraform plan -out=tfplan
```

### 4. Apply Configuration
```bash
terraform apply tfplan
```

### 5. Get Outputs
```bash
terraform output
```

## Environment Variables

### Backend Services
```
DATABASE_URL=postgresql://user:password@host:5432/db
REDIS_URL=redis://localhost:6379
KAFKA_BROKERS=localhost:9092
JWT_SECRET=your-jwt-secret
API_PORT=8080
LOG_LEVEL=info
```

### Mobile Apps
```
API_BASE_URL=https://api.famgo.com
GOOGLE_MAPS_API_KEY=your-key
FIREBASE_PROJECT_ID=your-project
```

## Database Migrations

### Run Migrations
```bash
psql -h localhost -U famgo -d famgo_platform -f database/migrations/001_initial_schema.sql
psql -h localhost -U famgo -d famgo_platform -f database/migrations/002_indexes.sql
# ... run all migrations in order
```

### Rollback
```bash
psql -h localhost -U famgo -d famgo_platform -f database/rollback/001_rollback.sql
```

## Monitoring & Observability

### Access Dashboards
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Jaeger**: http://localhost:16686

### View Logs
```bash
# Docker
docker-compose logs -f <service-name>

# Kubernetes
kubectl logs -f deployment/<service-name> -n famgo
```

## Testing

### Unit Tests
```bash
cd backend
go test ./... -v
```

### Integration Tests
```bash
go test -tags=integration ./tests/integration -v
```

### E2E Tests
```bash
go test -tags=e2e ./tests/e2e -v
```

## Security Checklist

- [ ] Change all default passwords
- [ ] Enable TLS/SSL certificates
- [ ] Configure firewall rules
- [ ] Set up VPN for internal communication
- [ ] Enable API rate limiting
- [ ] Configure CORS properly
- [ ] Encrypt sensitive data at rest
- [ ] Regular security audits
- [ ] Implement DDoS protection
- [ ] Set up backup strategy

## Performance Tuning

### Database
```sql
-- Connection pooling
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Index optimization
VACUUM ANALYZE;
REINDEX DATABASE famgo_platform;
```

### Redis
```
# Increase maxmemory
maxmemory 2gb
maxmemory-policy allkeys-lru
```

### Kafka
```
# Increase partitions
kafka-topics --alter --topic rides --partitions 10
```

## Troubleshooting

### Services Not Starting
```bash
# Check logs
docker-compose logs <service>

# Verify connectivity
curl http://localhost:8080/health
```

### Database Connection Issues
```bash
# Test connection
psql -h localhost -U famgo -d famgo_platform

# Check credentials
cat docker-compose.yml | grep POSTGRES
```

### High Latency
```bash
# Check metrics
curl http://localhost:9090/api/v1/query?query=http_request_duration

# Analyze slow queries
SELECT * FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;
```

## Backup & Recovery

### Database Backup
```bash
pg_dump -h localhost -U famgo -d famgo_platform > backup.sql
```

### Database Restore
```bash
psql -h localhost -U famgo -d famgo_platform < backup.sql
```

## Scaling

### Horizontal Scaling
```bash
# Kubernetes autoscaling
kubectl autoscale deployment ride-service --min=2 --max=10 -n famgo
```

### Vertical Scaling
```bash
# Increase resources in k8s
kubectl set resources deployment ride-service -c ride-service --limits=cpu=1,memory=1Gi
```

---

**Production Deployment Ready**: All components configured for enterprise-grade reliability and performance.
