package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	dob := time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC)
	at := time.Date(2025, 5, 9, 0, 0, 0, 0, time.UTC)
	if got := CalculateAge(dob, at); got != 34 {
		t.Fatalf("expected 34 got %d", got)
	}

	at2 := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
	if got := CalculateAge(dob, at2); got != 35 {
		t.Fatalf("expected 35 got %d", got)
	}
}
