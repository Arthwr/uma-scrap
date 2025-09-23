package main

import (
	"fmt"
	"log"

	"github.com/arthwr/uma-scrap/internal/config"
	"github.com/arthwr/uma-scrap/internal/scraper"
)

func main() {
	fmt.Printf("Starting parser at %s ...\n", config.BASE_URL)

	scr := scraper.NewScraper()
	scr.RegisterHandlers()

	if err := scr.Run(config.BASE_URL + config.EVENTS_URL); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Writing JSON file to %s ...\n", config.DEF_OUTPUT_DIR)
	if err := scr.Store().ExportJSON(config.DEF_OUTPUT_DIR); err != nil {
		log.Fatal("Failed to write JSON: ", err)
	}

	for key, count := range scr.Store().Counts {
		fmt.Printf("%s events: %d\n", key.String(), count)
	}

	fmt.Println("Job done!")
}
