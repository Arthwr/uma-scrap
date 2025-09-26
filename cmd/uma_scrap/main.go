package main

import (
	"log"

	"github.com/arthwr/uma-scrap/internal/service"
)

func main() {
	if err := service.RunScraper(); err != nil {
		log.Fatal(err)
	}
}
