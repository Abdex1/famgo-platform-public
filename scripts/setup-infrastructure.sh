#!/bin/bash
# ============================================================================
# FamGo Platform Infrastructure Setup Script
# ============================================================================
# This script automates the setup and verification of the entire infrastructure
#
# Usage: bash ./scripts/setup-infrastructure.sh [command]
# Commands: setup, verify, start, stop, clean, logs
# ============================================================================

set -e

# ============================================================================
# CONFIGURATION
# ============================================================================
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DOCKER_COMPOSE_FILE="$PROJECT_ROOT/infra/docker/docker-compose.yml"
ENV_FILE="$PROJECT_ROOT/.env.local"
LOG_DIR="$PROJECT_ROOT/logs"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ============================================================================
# UTILITY FUNCTIONS
# ============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi

    if [ ! -f "$ENV_FILE" ]; then
        log_error ".env.local file not found"
        log_info "Please create .env.local from .env.example"
        exit 1
    fi

    if [ ! -f "$DOCKER_COMPOSE_FILE" ]; then
        log_error "Docker Compose file not found: $DOCKER_COMPOSE_FILE"
        exit 1
    fi

    log_success "Prerequisites check passed"
}

create_log_directory() {
    if [ ! -d "$LOG_DIR" ]; then
        mkdir -p "$LOG_DIR"
        log_success "Created logs directory: $LOG_DIR"
    fi
}

# ============================================================================
# SETUP FUNCTION
# ============================================================================

setup() {
    log_info "Starting infrastructure setup..."

    check_prerequisites
    create_log_directory

    log_info "Building Docker images..."
    docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" build

    log_info "Creating volumes and networks..."
    docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" create

    log_success "Infrastructure setup completed!"
}

# ============================================================================
# START FUNCTION
# ============================================================================

start() {
    log_info "Starting containers..."

    docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" up -d

    log_info "Waiting for services to start (30 seconds)..."
    sleep 30

    log_success "Containers started!"
    log_info "Verifying services..."
    verify
}

# ============================================================================
# STOP FUNCTION
# ============================================================================

stop() {
    log_info "Stopping containers..."

    docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" down

    log_success "Containers stopped!"
}

# ============================================================================
# VERIFY FUNCTION
# ============================================================================

verify() {
    log_info "Verifying infrastructure..."

    local failed=0

    # PostgreSQL
    log_info "Checking PostgreSQL (port 5432)..."
    if docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" exec -T postgres pg_isready > /dev/null 2>&1; then
        log_success "PostgreSQL: OK"
    else
        log_error "PostgreSQL: FAILED"
        ((failed++))
    fi

    # Redis
    log_info "Checking Redis (port 6379)..."
    if docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" exec -T redis redis-cli ping | grep -q "PONG"; then
        log_success "Redis: OK"
    else
        log_error "Redis: FAILED"
        ((failed++))
    fi

    # Kafka
    log_info "Checking Kafka (port 9092)..."
    if docker ps | grep -q "famgo-kafka"; then
        log_success "Kafka: OK"
    else
        log_error "Kafka: FAILED"
        ((failed++))
    fi

    # MinIO
    log_info "Checking MinIO (port 9000)..."
    if curl -s http://localhost:9000/minio/health/live > /dev/null; then
        log_success "MinIO: OK"
    else
        log_error "MinIO: FAILED"
        ((failed++))
    fi

    # ClickHouse
    log_info "Checking ClickHouse (port 8123)..."
    if curl -s http://localhost:8123/ping > /dev/null; then
        log_success "ClickHouse: OK"
    else
        log_error "ClickHouse: FAILED"
        ((failed++))
    fi

    # Prometheus
    log_info "Checking Prometheus (port 9090)..."
    if curl -s http://localhost:9090/-/healthy > /dev/null; then
        log_success "Prometheus: OK"
    else
        log_error "Prometheus: FAILED"
        ((failed++))
    fi

    # Grafana
    log_info "Checking Grafana (port 3001)..."
    if curl -s http://localhost:3001/api/health > /dev/null; then
        log_success "Grafana: OK"
    else
        log_error "Grafana: FAILED"
        ((failed++))
    fi

    # Loki
    log_info "Checking Loki (port 3100)..."
    if curl -s http://localhost:3100/ready > /dev/null; then
        log_success "Loki: OK"
    else
        log_error "Loki: FAILED"
        ((failed++))
    fi

    # Jaeger
    log_info "Checking Jaeger (port 16686)..."
    if curl -s http://localhost:16686/api/services > /dev/null; then
        log_success "Jaeger: OK"
    else
        log_error "Jaeger: FAILED"
        ((failed++))
    fi

    if [ $failed -eq 0 ]; then
        log_success "All services verified successfully!"
        echo ""
        log_info "Access URLs:"
        echo "  PostgreSQL:    localhost:5432"
        echo "  Redis:         localhost:6379"
        echo "  Kafka:         localhost:9092"
        echo "  MinIO:         http://localhost:9000 (Console: http://localhost:9001)"
        echo "  ClickHouse:    http://localhost:8123"
        echo "  Prometheus:    http://localhost:9090"
        echo "  Grafana:       http://localhost:3001 (admin/admin_dev_password)"
        echo "  Loki:          http://localhost:3100"
        echo "  Jaeger:        http://localhost:16686"
        echo "  Nginx:         http://localhost:80"
        return 0
    else
        log_error "$failed service(s) failed verification"
        return 1
    fi
}

# ============================================================================
# CLEAN FUNCTION
# ============================================================================

clean() {
    log_warning "This will remove all containers, volumes, and data!"
    read -p "Are you sure? (yes/no): " -r
    echo
    if [[ $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        log_info "Cleaning up infrastructure..."

        docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" down -v

        log_success "Infrastructure cleaned!"
    else
        log_info "Cleanup cancelled"
    fi
}

# ============================================================================
# LOGS FUNCTION
# ============================================================================

logs() {
    local service=${1:-""}

    if [ -z "$service" ]; then
        docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" logs -f
    else
        docker-compose -f "$DOCKER_COMPOSE_FILE" --env-file "$ENV_FILE" logs -f "$service"
    fi
}

# ============================================================================
# MAIN
# ============================================================================

main() {
    local command="${1:-setup}"

    case "$command" in
        setup)
            setup
            ;;
        start)
            start
            ;;
        stop)
            stop
            ;;
        verify)
            verify
            ;;
        clean)
            clean
            ;;
        logs)
            logs "$2"
            ;;
        *)
            echo "Usage: $0 {setup|start|stop|verify|clean|logs [service]}"
            exit 1
            ;;
    esac
}

main "$@"
