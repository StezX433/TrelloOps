package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	TrelloAPIKey string
	TrelloToken  string

	TrelloBoardID string

	TrelloDoingListName string

	StaleDays int

	AtRiskLabelName  string
	AtRiskLabelColor string

	DryRun bool
}

func Load() (Config, error) {
	apiKey := os.Getenv("TRELLO_API_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	boardID := os.Getenv("TRELLO_BOARD_ID")

	if apiKey == "" {
		return Config{}, fmt.Errorf("TRELLO_API_KEY is missing")
	}

	if token == "" {
		return Config{}, fmt.Errorf("TRELLO_TOKEN is missing")
	}

	if boardID == "" {
		return Config{}, fmt.Errorf("TRELLO_BOARD_ID is missing")
	}

	doingListName := getEnvOrDefault(
		"TRELLO_DOING_LIST_NAME",
		"Doing",
	)

	staleDays, err := getIntEnv("STALE_DAYS", 5)
	if err != nil {
		return Config{}, err
	}

	if staleDays <= 0 {
		return Config{}, fmt.Errorf("STALE_DAYS must be greater than 0")
	}

	labelName := getEnvOrDefault(
		"AT_RISK_LABEL_NAME",
		"At Risk",
	)

	labelColor := getEnvOrDefault(
		"AT_RISK_LABEL_COLOR",
		"red",
	)

	dryRun, err := getBoolEnv("TRELLO_DRY_RUN", false)
	if err != nil {
		return Config{}, err
	}

	return Config{
		TrelloAPIKey:        apiKey,
		TrelloToken:         token,
		TrelloBoardID:       boardID,
		TrelloDoingListName: doingListName,
		StaleDays:           staleDays,
		AtRiskLabelName:     labelName,
		AtRiskLabelColor:    labelColor,
		DryRun:              dryRun,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getIntEnv(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf(
			"%s must be an integer: %w",
			key,
			err,
		)
	}

	return result, nil
}

func getBoolEnv(key string, defaultValue bool) (bool, error) {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue, nil
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf(
			"%s must be true or false: %w",
			key,
			err,
		)
	}

	return result, nil
}
