/*
# PHASE 5 — REFRESH TOKEN ROTATION

# CREATE TOKEN HASHER

internal/infrastructure/security/token_hasher.go
*/
package security

import (
    "crypto/sha256"
    "encoding/hex"
)

func HashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}
