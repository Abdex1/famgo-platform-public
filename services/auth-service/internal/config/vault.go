
/*# PHASE 11 — AUTH-SERVICE VAULT INTEGRATION

Now integrate Vault into auth-service.

---

# STEP 17 — CREATE VAULT CONFIG

services/auth-service/internal/config/vault.go
*/
package config

type VaultConfig struct {
    Enabled bool

    Address string
    Token   string

    MountPath string
    KVVersion int
}
