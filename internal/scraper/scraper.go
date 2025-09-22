package scraper

import (
	"fmt"
	"log"
	"time"

	"github.com/arthwr/uma-scrap/internal/config"
	"github.com/arthwr/uma-scrap/internal/models"
	"github.com/gocolly/colly"
)

type Scraper struct {
	collector *colly.Collector
	store     *models.EventStore
}

func NewScraper() *Scraper {
	c := colly.NewCollector(
		colly.AllowedDomains(config.Domain),
		colly.MaxDepth(config.MaxDepth),
		colly.Async(config.Async),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  config.Glob,
		Parallelism: config.Workers,
		Delay:       config.RequestDelay * time.Second,
		RandomDelay: config.RequestRandomDelay * time.Second,
	})

	return &Scraper{
		collector: c,
		store:     models.NewEventStore(),
	}
}

func (s *Scraper) RegisterHandlers() {
	s.collector.OnHTML(SelectorEventsTable, func(el *colly.HTMLElement) {
		el.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
			event, err := parseEventRow(tr)
			if err != nil {
				log.Println("Skipping invalid row: ", err)
				return
			}

			s.store.AddEvent(event)
		})
	})

	s.collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})
}

func (s *Scraper) Run(startURL string) error {
	if err := s.collector.Visit(startURL); err != nil {
		return err
	}
	s.collector.Wait()
	return nil
}

func (s *Scraper) Store() *models.EventStore {
	return s.store
}
