package rules

import (
	"testing"
	"time"

	"trelloops/internal/trello"
)

func TestIsStale(t *testing.T) {
	now := time.Date(
		2026,
		10,
		6,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	threshold := 5 * 24 * time.Hour

	tests := []struct {
		name          string
		lastActivity  string
		expectedStale bool
	}{
		{
			name:          "recent card",
			lastActivity:  "2026-10-04T12:00:00Z",
			expectedStale: false,
		},
		{
			name:          "exactly five days old",
			lastActivity:  "2026-10-01T12:00:00Z",
			expectedStale: true,
		},
		{
			name:          "older than five days",
			lastActivity:  "2026-09-25T12:00:00Z",
			expectedStale: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := trello.Card{
				Name:             "Test Card",
				DateLastActivity: tt.lastActivity,
			}

			result, err := IsStale(
				card,
				now,
				threshold,
			)

			if err != nil {
				t.Fatal(err)
			}

			if result != tt.expectedStale {
				t.Fatalf(
					"expected stale=%v, got %v",
					tt.expectedStale,
					result,
				)
			}
		})
	}
}

func TestIsStaleInvalidDate(t *testing.T) {
	card := trello.Card{
		Name:             "Broken Card",
		DateLastActivity: "invalid-date",
	}

	_, err := IsStale(
		card,
		time.Now(),
		5*24*time.Hour,
	)

	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}
