FROM kong:3.0-alpine

# Install additional tools
RUN apk add --no-cache curl

# Copy Kong configuration
COPY kong.yml /etc/kong/kong.yml
COPY kong-init.sh /docker-entrypoint.d/

# Make init script executable
RUN chmod +x /docker-entrypoint.d/kong-init.sh

# Expose ports
EXPOSE 8000 8443 8001 8444

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8001/_healthcheck || exit 1

# Run Kong
CMD ["kong", "start"]
