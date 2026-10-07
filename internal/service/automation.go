package service

import (
	"fmt"
	"time"

	"trelloops/internal/rules"
	"trelloops/internal/trello"
)

type TrelloAPI interface {
	GetLists(boardID string) ([]trello.List, error)
	GetCards(listID string) ([]trello.Card, error)

	GetLabels(boardID string) ([]trello.Label, error)

	CreateLabel(
		boardID string,
		name string,
		color string,
	) (*trello.Label, error)

	AddLabelToCard(
		cardID string,
		labelID string,
	) error
}

type Automation struct {
	client TrelloAPI

	boardID string

	doingListName string

	staleThreshold time.Duration

	labelName  string
	labelColor string

	dryRun bool

	now func() time.Time
}

type Result struct {
	CardsChecked int
	StaleCards   int
	LabeledCards int
	SkippedCards int
}

func NewAutomation(
	client TrelloAPI,
	boardID string,
	doingListName string,
	staleThreshold time.Duration,
	labelName string,
	labelColor string,
	dryRun bool,
) *Automation {
	return &Automation{
		client:         client,
		boardID:        boardID,
		doingListName:  doingListName,
		staleThreshold: staleThreshold,
		labelName:      labelName,
		labelColor:     labelColor,
		dryRun:         dryRun,
		now:            time.Now,
	}
}

func (a *Automation) Run() (Result, error) {
	var result Result

	lists, err := a.client.GetLists(a.boardID)
	if err != nil {
		return result, err
	}

	doingListID := ""

	for _, list := range lists {
		if list.Name == a.doingListName {
			doingListID = list.ID
			break
		}
	}

	if doingListID == "" {
		return result, fmt.Errorf(
			"list %q not found",
			a.doingListName,
		)
	}

	cards, err := a.client.GetCards(doingListID)
	if err != nil {
		return result, err
	}

	label, err := a.findOrCreateLabel()
	if err != nil {
		return result, err
	}

	now := a.now()

	for _, card := range cards {
		result.CardsChecked++

		stale, err := rules.IsStale(
			card,
			now,
			a.staleThreshold,
		)

		if err != nil {
			return result, err
		}

		if !stale {
			continue
		}

		result.StaleCards++

		if hasLabel(card, label.ID) {
			result.SkippedCards++
			continue
		}

		if a.dryRun {
			result.LabeledCards++
			continue
		}

		if err := a.client.AddLabelToCard(
			card.ID,
			label.ID,
		); err != nil {
			return result, fmt.Errorf(
				"failed to label card %q: %w",
				card.Name,
				err,
			)
		}

		result.LabeledCards++
	}

	return result, nil
}

func (a *Automation) findOrCreateLabel() (*trello.Label, error) {
	labels, err := a.client.GetLabels(a.boardID)
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		if label.Name == a.labelName {
			return &label, nil
		}
	}

	if a.dryRun {
		return &trello.Label{
			ID:      "dry-run-label",
			Name:    a.labelName,
			Color:   a.labelColor,
			IDBoard: a.boardID,
		}, nil
	}

	return a.client.CreateLabel(
		a.boardID,
		a.labelName,
		a.labelColor,
	)
}

func hasLabel(card trello.Card, labelID string) bool {
	for _, id := range card.IDLabels {
		if id == labelID {
			return true
		}
	}

	return false
}
