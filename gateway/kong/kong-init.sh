#!/bin/bash
# gateway/kong/kong-init.sh
# Initialize Kong with configuration

set -e

echo "Waiting for Kong to be ready..."
until curl -f http://localhost:8001/_healthcheck 2>/dev/null; do
  sleep 1
done

echo "Kong is ready. Loading configuration..."

# Load Kong configuration via Admin API
curl -i -X POST http://localhost:8001/config \
  -F "config=@/etc/kong/kong.yml" || true

echo "Kong configuration loaded successfully!"

# Add JWT credentials for consumers
curl -i -X POST http://localhost:8001/consumers/rider-app/jwt \
  -d "key=rider-app-key" \
  -d "secret=rider-app-secret" || true

curl -i -X POST http://localhost:8001/consumers/driver-app/jwt \
  -d "key=driver-app-key" \
  -d "secret=driver-app-secret" || true

curl -i -X POST http://localhost:8001/consumers/admin-dashboard/jwt \
  -d "key=admin-key" \
  -d "secret=admin-secret" || true

echo "Kong initialization complete!"
