#!/bin/bash

# FamGo Platform - Kafka Topics Setup Script
# Creates all 15 Kafka topics for event-driven architecture

set -e

KAFKA_BROKERS="${KAFKA_BROKERS:-localhost:9092}"
KAFKA_CMD="kafka-topics.sh --bootstrap-server $KAFKA_BROKERS"

echo "=== FamGo Kafka Topics Setup ==="
echo "Kafka Brokers: $KAFKA_BROKERS"
echo ""

# Function to create topic
create_topic() {
    local topic=$1
    echo "Creating topic: $topic"
    
    $KAFKA_CMD --create \
        --topic $topic \
        --partitions 3 \
        --replication-factor 1 \
        --config retention.ms=604800000 \
        --config compression.type=snappy \
        --if-not-exists
}

# Function to create consumer group
create_consumer_group() {
    local group=$1
    local topics=$2
    
    echo "Creating consumer group: $group"
    
    kafka-consumer-groups.sh --bootstrap-server $KAFKA_BROKERS \
        --create \
        --group $group \
        --reset-offsets --to-earliest \
        --topic $topics \
        --execute || true
}

echo "--- Creating Topics ---"

# Ride Events
create_topic "ride.created"
create_topic "ride.matching.started"
create_topic "ride.driver.assigned"
create_topic "ride.started"
create_topic "ride.completed"
create_topic "ride.cancelled"

# Location Events
create_topic "driver.location.updated"

# Pooling Events
create_topic "pool.created"
create_topic "pool.updated"

# Pricing Events
create_topic "pricing.calculated"

# Payment Events
create_topic "payment.completed"
create_topic "payment.failed"

# Wallet Events
create_topic "wallet.transaction.created"

# Safety Events
create_topic "safety.sos.triggered"

# Fraud Events
create_topic "fraud.detected"

# Notification Events
create_topic "notification.send"

echo ""
echo "--- Creating Consumer Groups ---"

# Service consumer groups
create_consumer_group "auth-service-group" "ride.*,payment.*"
create_consumer_group "ride-service-group" "pool.*,payment.*,driver.*"
create_consumer_group "dispatch-service-group" "ride.created"
create_consumer_group "pooling-service-group" "ride.created,driver.location.updated"
create_consumer_group "gps-service-group" "ride.started,ride.completed"
create_consumer_group "payment-service-group" "ride.completed"
create_consumer_group "wallet-service-group" "payment.completed,payment.failed"
create_consumer_group "notification-service-group" "ride.*,payment.*,safety.*,fraud.*"
create_consumer_group "safety-service-group" "ride.started,ride.completed"
create_consumer_group "fraud-service-group" "payment.completed,ride.created"
create_consumer_group "pricing-service-group" "ride.created,ride.started"
create_consumer_group "analytics-service-group" "ride.*,payment.*,driver.*,pool.*"

echo ""
echo "--- Verifying Topics ---"

$KAFKA_CMD --list

echo ""
echo "=== Kafka Setup Complete ==="
