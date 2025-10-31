package utils

import (
    "testing"
    "time"
)

func TestConvertTimeToDBDateTime(t *testing.T) {
    ts := time.Date(2024, time.January, 2, 15, 4, 5, 0, time.UTC)
    got := ConvertTimeToDBDateTime(ts)
    if got != "2024-01-02 15:04:05" {
        t.Fatalf("expected formatted time, got %s", got)
    }
}

func TestGetCurrentDateTime(t *testing.T) {
    before := time.Now()
    current := GetCurrentDateTime()
    after := time.Now()
    if current.Before(before) || current.After(after.Add(time.Second)) {
        t.Fatalf("GetCurrentDateTime returned unexpected time: %v", current)
    }
}

