"""
# PHASE 12 — VAULT POLICIES

Now we isolate services.

---

# STEP 20 — CREATE AUTH POLICY

security/vault/policies/auth-service-policy.hcl
"""
path "kv/data/auth-service" {
  capabilities = ["read"]
}