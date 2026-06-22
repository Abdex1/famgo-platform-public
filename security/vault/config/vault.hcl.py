"""
# PHASE 2 — ENTERPRISE VAULT DEV CONFIG

---

# STEP 2 — CREATE VAULT CONFIG


security/vault/config/vault.hcl
"""
ui = true

disable_mlock = true

api_addr = "http://0.0.0.0:8200"
cluster_addr = "http://0.0.0.0:8201"

storage "file" {
  path = "/vault/data"
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
}

telemetry {
  prometheus_retention_time = "24h"
  disable_hostname = true
}