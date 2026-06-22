#!/bin/bash

# FamGo Platform - Redis Setup Script
# Initializes Redis cache, GEO indexes, and session storage

set -e

REDIS_HOST="${REDIS_HOST:-localhost}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_DB="${REDIS_DB:-0}"

echo "=== FamGo Redis Setup ==="
echo "Redis Host: $REDIS_HOST"
echo "Redis Port: $REDIS_PORT"

# Wait for Redis to be ready
echo "Waiting for Redis to be ready..."
for i in {1..30}; do
  if redis-cli -h $REDIS_HOST -p $REDIS_PORT ping >/dev/null 2>&1; then
    echo "✓ Redis is ready"
    break
  fi
  echo "Attempt $i/30..."
  sleep 1
done

# Create Redis CLI commands
redis-cli -h $REDIS_HOST -p $REDIS_PORT << EOF

# ============================================================================
# SELECT DATABASE
# ============================================================================
SELECT $REDIS_DB

# ============================================================================
# GEOSPATIAL INDEX - ACTIVE DRIVERS
# ============================================================================

# Initialize drivers GEO index (longitude, latitude, member)
# This will be updated in real-time by GPS service

# Example data (Addis Ababa locations):
GEOADD drivers:geo 9.0320 8.9868 "driver:uuid-1"
GEOADD drivers:geo 9.0355 8.9867 "driver:uuid-2"
GEOADD drivers:geo 9.0210 8.9756 "driver:uuid-3"

# ============================================================================
# SESSION STORAGE
# ============================================================================

# Session hash structure:
# HSET session:{session_id} user_id {user_id} device_id {device_id} ...
# EXPIRE session:{session_id} 3600  # 1 hour TTL

# Example:
# HSET session:abc123 user_id user-1 role rider is_active true

# ============================================================================
# RATE LIMITING - Sliding Window Counters
# ============================================================================

# Rate limit keys (auto-expire with TTL):
# INCR rate:user:{user_id}:requests:{minute}
# EXPIRE rate:user:{user_id}:requests:{minute} 60

# ============================================================================
# OTP STORAGE
# ============================================================================

# OTP codes (with TTL):
# SET otp:{phone_number}:{purpose} {code} EX 600  # 10 minutes
# Example:
# SET otp:+251911234567:login ABC123 EX 600

# ============================================================================
# CACHE LAYERS
# ============================================================================

# User profiles (cache):
# HSET cache:user:{user_id} email name phone role
# EXPIRE cache:user:{user_id} 3600  # 1 hour

# Driver profiles (cache):
# HSET cache:driver:{driver_id} rating acceptance_rate status
# EXPIRE cache:driver:{driver_id} 1800  # 30 minutes

# ============================================================================
# REALTIME PRESENCE
# ============================================================================

# Driver online status:
# SET driver:online:{driver_id} true EX 300  # 5 minute heartbeat
# Example:
# SET driver:online:driver-1 true EX 300

# ============================================================================
# LOCKS & DISTRIBUTED TRANSACTIONS
# ============================================================================

# Distributed locks (for critical operations):
# SET lock:{resource_id} {lock_token} EX 30 NX
# Example:
# SET lock:ride:ride-1:processing abc123 EX 30 NX

# ============================================================================
# ANALYTICS COUNTERS (Real-time metrics)
# ============================================================================

# Daily ride counter:
# INCR metrics:rides:count:{date}
# EXPIRE metrics:rides:count:{date} 86400

# Hourly earnings:
# INCR metrics:earnings:{hour}
# EXPIRE metrics:earnings:{hour} 3600

# ============================================================================
# CONFIGURATION
# ============================================================================

# Set max memory policy
CONFIG SET maxmemory-policy allkeys-lru
CONFIG SET maxmemory 1gb

# Enable keyspace notifications (for expiration events)
CONFIG SET notify-keyspace-events "Ex"

# ============================================================================
# INFO & MONITORING
# ============================================================================

# Get Redis info
INFO server
INFO memory
INFO stats
INFO replication

# ============================================================================
# PERSISTENCE
# ============================================================================

# Enable AOF (Append-Only File)
CONFIG SET appendonly yes
CONFIG SET appendfsync everysec

# Enable RDB snapshots
CONFIG SET save "900 1 300 10 60 10000"

EOF

echo ""
echo "=== Redis Setup Complete ==="
echo ""
echo "Redis Keyspaces:"
echo "  - drivers:geo              → GEO index of active drivers"
echo "  - session:{id}             → User sessions"
echo "  - rate:user:{id}:*         → Rate limiting counters"
echo "  - otp:{phone}:{purpose}    → OTP codes"
echo "  - cache:user:{id}          → User profile cache"
echo "  - cache:driver:{id}        → Driver profile cache"
echo "  - driver:online:{id}       → Driver presence"
echo "  - lock:{resource}          → Distributed locks"
echo "  - metrics:*                → Real-time analytics"
echo ""
echo "Common Operations:"
echo "  # Find nearby drivers (within 5km)"
echo "  GEORADIUSBYMEMBER drivers:geo driver:uuid-1 5 km"
echo ""
echo "  # Get driver position"
echo "  GEOPOS drivers:geo driver:uuid-1"
echo ""
echo "  # Distance between drivers"
echo "  GEODIST drivers:geo driver:uuid-1 driver:uuid-2 km"
echo ""
echo "  # Set OTP"
echo "  SET otp:+251911234567:login 123456 EX 600"
echo ""
echo "  # Verify OTP"
echo "  GET otp:+251911234567:login"
echo ""
