#!/bin/bash
# Kafka Topic Setup Script for FamGo Platform
# ============================================================================
# This script creates all Kafka topics required for the platform
# Run this after Kafka is started: docker exec kafka ./topics-setup.sh

set -e

KAFKA_BROKER="${KAFKA_BROKER:-kafka:9092}"
REPLICATION_FACTOR="${REPLICATION_FACTOR:-1}"
PARTITIONS="${PARTITIONS:-3}"

echo "Creating Kafka topics for FamGo Platform..."
echo "Broker: $KAFKA_BROKER"
echo "Replication Factor: $REPLICATION_FACTOR"
echo "Partitions: $PARTITIONS"

# Function to create topic
create_topic() {
    local topic_name=$1
    local num_partitions=${2:-$PARTITIONS}
    local replication=${3:-$REPLICATION_FACTOR}
    
    echo "Creating topic: $topic_name"
    kafka-topics --create \
        --bootstrap-server "$KAFKA_BROKER" \
        --topic "$topic_name" \
        --partitions "$num_partitions" \
        --replication-factor "$replication" \
        --if-not-exists || echo "Topic $topic_name already exists"
}

# ============================================================================
# RIDE MANAGEMENT TOPICS
# ============================================================================
create_topic "ride.created" 3
create_topic "ride.matching.started" 3
create_topic "ride.driver.assigned" 3
create_topic "ride.started" 3
create_topic "ride.in_progress" 3
create_topic "ride.completed" 3
create_topic "ride.cancelled" 3
create_topic "ride.updated" 3

# ============================================================================
# DRIVER MANAGEMENT TOPICS
# ============================================================================
create_topic "driver.location.updated" 6 1  # High partition count for GPS
create_topic "driver.status.changed" 3
create_topic "driver.online" 3
create_topic "driver.offline" 3
create_topic "driver.available" 3
create_topic "driver.busy" 3

# ============================================================================
# POOL MANAGEMENT TOPICS
# ============================================================================
create_topic "pool.created" 3
create_topic "pool.updated" 3
create_topic "pool.passenger.added" 3
create_topic "pool.passenger.removed" 3
create_topic "pool.closed" 3

# ============================================================================
# PRICING TOPICS
# ============================================================================
create_topic "pricing.calculated" 3
create_topic "pricing.surge.triggered" 3
create_topic "pricing.discount.applied" 3
create_topic "pricing.updated" 3

# ============================================================================
# PAYMENT TOPICS
# ============================================================================
create_topic "payment.initiated" 3
create_topic "payment.completed" 3
create_topic "payment.failed" 3
create_topic "payment.refunded" 3
create_topic "payment.verified" 3

# ============================================================================
# WALLET TOPICS
# ============================================================================
create_topic "wallet.transaction.created" 3
create_topic "wallet.transaction.completed" 3
create_topic "wallet.transaction.failed" 3
create_topic "wallet.balance.updated" 3
create_topic "wallet.topup.requested" 3

# ============================================================================
# SAFETY TOPICS
# ============================================================================
create_topic "safety.sos.triggered" 3
create_topic "safety.panic.alert" 3
create_topic "safety.anomaly.detected" 3
create_topic "safety.incident.reported" 3
create_topic "safety.verification.requested" 3

# ============================================================================
# FRAUD DETECTION TOPICS
# ============================================================================
create_topic "fraud.detected" 3
create_topic "fraud.score.calculated" 3
create_topic "fraud.alert.triggered" 3
create_topic "fraud.review.requested" 3

# ============================================================================
# NOTIFICATION TOPICS
# ============================================================================
create_topic "notification.send" 3
create_topic "notification.email.sent" 3
create_topic "notification.sms.sent" 3
create_topic "notification.push.sent" 3
create_topic "notification.sent.failed" 3

# ============================================================================
# SUBSCRIPTION TOPICS
# ============================================================================
create_topic "subscription.created" 3
create_topic "subscription.activated" 3
create_topic "subscription.cancelled" 3
create_topic "subscription.renewed" 3

# ============================================================================
# ANALYTICS & REPORTING TOPICS
# ============================================================================
create_topic "analytics.ride.metrics" 3
create_topic "analytics.driver.metrics" 3
create_topic "analytics.user.metrics" 3
create_topic "analytics.payment.metrics" 3
create_topic "analytics.operational.metrics" 3

# ============================================================================
# SYSTEM TOPICS
# ============================================================================
create_topic "system.health.check" 3
create_topic "system.error" 3
create_topic "system.warning" 3
create_topic "system.audit.log" 3
create_topic "system.config.change" 3

# ============================================================================
# LIST ALL TOPICS
# ============================================================================
echo ""
echo "Kafka topics created successfully!"
echo "Listing all topics:"
kafka-topics --list --bootstrap-server "$KAFKA_BROKER"

echo "Setup complete!"
