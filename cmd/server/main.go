package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"trelloops/internal/config"
	"trelloops/internal/service"
	"trelloops/internal/trello"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf(
				"failed to load .env: %v",
				err,
			)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client := trello.NewClient(
		cfg.TrelloAPIKey,
		cfg.TrelloToken,
	)

	automation := service.NewAutomation(
		client,
		cfg.TrelloBoardID,
		cfg.TrelloDoingListName,
		time.Duration(cfg.StaleDays)*24*time.Hour,
		cfg.AtRiskLabelName,
		cfg.AtRiskLabelColor,
		cfg.DryRun,
	)

	fmt.Println("===================================")
	fmt.Println("          TrelloOps")
	fmt.Println("===================================")
	fmt.Println()

	if cfg.DryRun {
		fmt.Println("Mode: DRY RUN")
	} else {
		fmt.Println("Mode: LIVE")
	}

	fmt.Println()
	fmt.Printf(
		"Board: %s\n",
		cfg.TrelloBoardID,
	)
	fmt.Printf(
		"List: %s\n",
		cfg.TrelloDoingListName,
	)
	fmt.Printf(
		"Stale threshold: %d days\n",
		cfg.StaleDays,
	)
	fmt.Printf(
		"Label: %s\n",
		cfg.AtRiskLabelName,
	)
	fmt.Println()

	result, err := automation.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-----------------------------------")
	fmt.Println("Automation complete")
	fmt.Println("-----------------------------------")

	fmt.Printf(
		"Cards checked: %d\n",
		result.CardsChecked,
	)

	fmt.Printf(
		"Stale cards: %d\n",
		result.StaleCards,
	)

	fmt.Printf(
		"Cards labeled: %d\n",
		result.LabeledCards,
	)

	fmt.Printf(
		"Cards skipped: %d\n",
		result.SkippedCards,
	)
}
