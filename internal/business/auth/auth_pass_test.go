package bussinessAuth

import "testing"

func TestHashAndComparePasswordArgon2id(t *testing.T) {
    password := "super-secret"

    hashed, err := HashPasswordArgon2id(password)
    if err != nil {
        t.Fatalf("expected hash generation to succeed, got error %v", err)
    }

    if hashed == "" {
        t.Fatalf("expected non-empty hash")
    }

    ok, err := ComparePasswordArgon2id(hashed, password)
    if err != nil {
        t.Fatalf("expected comparison to succeed, got error %v", err)
    }

    if !ok {
        t.Fatalf("expected password comparison to succeed")
    }

    ok, err = ComparePasswordArgon2id(hashed, "wrong-password")
    if err != nil {
        t.Fatalf("expected comparison with wrong password to return false without error, got %v", err)
    }

    if ok {
        t.Fatalf("expected password comparison to fail for wrong password")
    }

    if _, err := ComparePasswordArgon2id("invalid-hash", password); err == nil {
        t.Fatalf("expected error for invalid hash format")
    }
}

