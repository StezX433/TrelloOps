package rules

import (
	"fmt"
	"time"

	"trelloops/internal/trello"
)

func IsStale(card trello.Card, now time.Time, threshold time.Duration) (bool, error) {
	if card.DateLastActivity == "" {
		return false, fmt.Errorf(
			"card %q has no dateLastActivity",
			card.Name,
		)
	}

	lastActivity, err := time.Parse(
		time.RFC3339,
		card.DateLastActivity,
	)
	if err != nil {
		return false, fmt.Errorf(
			"invalid dateLastActivity for card %q: %w",
			card.Name,
			err,
		)
	}

	inactiveFor := now.Sub(lastActivity)

	return inactiveFor >= threshold, nil
}
