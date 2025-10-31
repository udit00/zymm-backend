package bussinessAuth

import (
    "os"
    "testing"
    "time"
)

func withSecretKey(t *testing.T, key string, fn func()) {
    t.Helper()
    oldEnv := os.Getenv("ZYMM_BACKEND_AUTH_SECRET_KEY")
    if err := os.Setenv("ZYMM_BACKEND_AUTH_SECRET_KEY", key); err != nil {
        t.Fatalf("failed to set env: %v", err)
    }
    MY_SECRET_KEY = ""
    defer func() {
        if oldEnv == "" {
            os.Unsetenv("ZYMM_BACKEND_AUTH_SECRET_KEY")
        } else {
            os.Setenv("ZYMM_BACKEND_AUTH_SECRET_KEY", oldEnv)
        }
        MY_SECRET_KEY = ""
    }()
    fn()
}

func TestGenerateAndValidateJWTToken(t *testing.T) {
    withSecretKey(t, "test-secret", func() {
        token, err := GenerateJWTToken(42, 3)
        if err != nil {
            t.Fatalf("expected token generation to succeed, got %v", err)
        }
        if token == nil || *token == "" {
            t.Fatalf("expected non-empty token string")
        }

        parsed, err := ValidateToken(*token)
        if err != nil {
            t.Fatalf("expected validation to succeed, got %v", err)
        }
        if !parsed.Valid {
            t.Fatalf("expected parsed token to be valid")
        }

        claims, err := GetDataFromJWTToken(*token)
        if err != nil {
            t.Fatalf("expected claims extraction to succeed, got %v", err)
        }
        if claims.UserId != 42 || claims.RoleId != 3 {
            t.Fatalf("unexpected claims data: %+v", claims)
        }
        if time.Unix(claims.Expiry, 0).Before(time.Now()) {
            t.Fatalf("expected token expiry to be in the future")
        }
    })
}

func TestValidateTokenWithWrongSecret(t *testing.T) {
    withSecretKey(t, "secret-one", func() {
        token, err := GenerateJWTToken(1, 1)
        if err != nil {
            t.Fatalf("expected token generation, got %v", err)
        }

        withSecretKey(t, "secret-two", func() {
            if _, err := ValidateToken(*token); err == nil {
                t.Fatalf("expected validation to fail with different secret")
            }
        })
    })
}

func TestGetDataFromJWTTokenInvalid(t *testing.T) {
    withSecretKey(t, "secret", func() {
        if _, err := GetDataFromJWTToken("invalid-token"); err == nil {
            t.Fatalf("expected error when parsing invalid token string")
        }
    })
}

