package main

import (
	"fmt"

	"github.com/arthwr/uma-scrap/internal/config"
	"github.com/arthwr/uma-scrap/internal/scraper"
)

func main() {
	fmt.Printf("Starting parser at %s\n", config.BaseURL)

	scr := scraper.NewScraper()
	scr.RegisterHandlers()

	if err := scr.Run(config.BaseURL + config.EventsURL); err != nil {
		fmt.Println("Scraper failed: ", err)
	}
}
