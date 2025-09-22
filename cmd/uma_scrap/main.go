package main

import (
	"fmt"
	"log"

	"github.com/arthwr/uma-scrap/internal/config"
	"github.com/arthwr/uma-scrap/internal/scraper"
)

func main() {
	fmt.Printf("Starting parser at %s ...\n", config.BaseURL)

	scr := scraper.NewScraper()
	scr.RegisterHandlers()

	if err := scr.Run(config.BaseURL + config.EventsURL); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Writting JSON file to %s ...\n", config.DefEventsFilename)
	if err := scr.Store().ExportJSON(config.DefOutputDir, config.DefEventsFilename); err != nil {
		log.Fatal("Failed to write JSON: ", err)
	}
	fmt.Println("Job done!")
}
