package utils

import "testing"

func TestConvertStringToInt(t *testing.T) {
    if val := ConvertStringToInt("42"); val == nil || *val != 42 {
        t.Fatalf("expected pointer to 42, got %v", val)
    }

    if val := ConvertStringToInt("invalid"); val != nil {
        t.Fatalf("expected nil for invalid input, got %v", *val)
    }
}

func TestConvertIntToString(t *testing.T) {
    if out := ConvertIntToString(55); out != "55" {
        t.Fatalf("expected '55', got %s", out)
    }
}

