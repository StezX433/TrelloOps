package service

import (
	"testing"
	"time"

	"trelloops/internal/trello"
)

type fakeTrelloClient struct {
	lists  []trello.List
	cards  []trello.Card
	labels []trello.Label

	createdLabels int
	addedLabels   []string
}

func (f *fakeTrelloClient) GetLists(boardID string) ([]trello.List, error) {
	return f.lists, nil
}

func (f *fakeTrelloClient) GetCards(listID string) ([]trello.Card, error) {
	return f.cards, nil
}

func (f *fakeTrelloClient) GetLabels(boardID string) ([]trello.Label, error) {
	return f.labels, nil
}

func (f *fakeTrelloClient) CreateLabel(
	boardID string,
	name string,
	color string,
) (*trello.Label, error) {
	f.createdLabels++

	label := trello.Label{
		ID:      "label-1",
		Name:    name,
		Color:   color,
		IDBoard: boardID,
	}

	f.labels = append(f.labels, label)

	return &label, nil
}

func (f *fakeTrelloClient) AddLabelToCard(
	cardID string,
	labelID string,
) error {
	f.addedLabels = append(
		f.addedLabels,
		cardID+":"+labelID,
	)

	return nil
}

func TestAutomationLabelsStaleCards(t *testing.T) {
	client := &fakeTrelloClient{
		lists: []trello.List{
			{
				ID:   "doing-1",
				Name: "Doing",
			},
		},
		cards: []trello.Card{
			{
				ID:               "stale-1",
				Name:             "Stale Card",
				IDList:           "doing-1",
				DateLastActivity: "2026-09-30T12:00:00Z",
			},
			{
				ID:               "active-1",
				Name:             "Active Card",
				IDList:           "doing-1",
				DateLastActivity: "2026-10-05T12:00:00Z",
			},
		},
	}

	automation := NewAutomation(
		client,
		"board-1",
		"Doing",
		5*24*time.Hour,
		"At Risk",
		"red",
		false,
	)

	automation.now = func() time.Time {
		return time.Date(
			2026,
			10,
			6,
			12,
			0,
			0,
			0,
			time.UTC,
		)
	}

	result, err := automation.Run()

	if err != nil {
		t.Fatal(err)
	}

	if result.CardsChecked != 2 {
		t.Fatalf(
			"expected 2 cards checked, got %d",
			result.CardsChecked,
		)
	}

	if result.StaleCards != 1 {
		t.Fatalf(
			"expected 1 stale card, got %d",
			result.StaleCards,
		)
	}

	if result.LabeledCards != 1 {
		t.Fatalf(
			"expected 1 labeled card, got %d",
			result.LabeledCards,
		)
	}

	if len(client.addedLabels) != 1 {
		t.Fatalf(
			"expected 1 label operation, got %d",
			len(client.addedLabels),
		)
	}

	if client.addedLabels[0] != "stale-1:label-1" {
		t.Fatalf(
			"unexpected label operation: %s",
			client.addedLabels[0],
		)
	}
}

func TestAutomationDoesNotDuplicateLabel(t *testing.T) {
	client := &fakeTrelloClient{
		lists: []trello.List{
			{
				ID:   "doing-1",
				Name: "Doing",
			},
		},
		cards: []trello.Card{
			{
				ID:               "stale-1",
				Name:             "Already At Risk",
				IDList:           "doing-1",
				DateLastActivity: "2026-09-20T12:00:00Z",
				IDLabels:         []string{"label-1"},
			},
		},
		labels: []trello.Label{
			{
				ID:      "label-1",
				Name:    "At Risk",
				Color:   "red",
				IDBoard: "board-1",
			},
		},
	}

	automation := NewAutomation(
		client,
		"board-1",
		"Doing",
		5*24*time.Hour,
		"At Risk",
		"red",
		false,
	)

	automation.now = func() time.Time {
		return time.Date(
			2026,
			10,
			6,
			12,
			0,
			0,
			0,
			time.UTC,
		)
	}

	result, err := automation.Run()

	if err != nil {
		t.Fatal(err)
	}

	if result.StaleCards != 1 {
		t.Fatalf(
			"expected 1 stale card, got %d",
			result.StaleCards,
		)
	}

	if result.LabeledCards != 0 {
		t.Fatalf(
			"expected 0 new labels, got %d",
			result.LabeledCards,
		)
	}

	if result.SkippedCards != 1 {
		t.Fatalf(
			"expected 1 skipped card, got %d",
			result.SkippedCards,
		)
	}

	if len(client.addedLabels) != 0 {
		t.Fatalf(
			"expected no label API calls",
		)
	}
}
