package service

import (
	"fmt"
	"log"

	"github.com/arthwr/uma-scrap/internal/config"
	"github.com/arthwr/uma-scrap/internal/scraper"
)

func RunScraper() error {
	fmt.Printf("Starting parser at %s ...\n", config.BASE_URL)

	scr := scraper.NewScraper()

	targetURL := config.BASE_URL + config.EVENTS_URL

	if err := scr.ScrapeEventList(targetURL); err != nil {
		return fmt.Errorf("failed scraping event list: %w", err)
	}

	if err := scr.ScrapeEventDetails(); err != nil {
		return fmt.Errorf("failed scraping event details: %w", err)
	}

	fmt.Printf("Writing JSON file to %s ...\n", config.DEF_OUTPUT_DIR)
	if err := scr.Store().ExportJSON(config.DEF_OUTPUT_DIR); err != nil {
		log.Fatal("failed to write JSON: ", err)
	}

	fmt.Println("\nResults:")
	for key, count := range scr.Store().Counts {
		fmt.Printf("	%s events: %d\n", key.String(), count)
	}

	fmt.Println("Job done!")
	return nil
}
