package scraper

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/arthwr/uma-scrap/internal/config"
	"github.com/arthwr/uma-scrap/internal/models"
	"github.com/gocolly/colly"
)

type EventLink struct {
	Name string
	URL  string
}

type Scraper struct {
	jobMutex      sync.Mutex
	collector     *colly.Collector
	store         *models.EventStore
	eventURLs     []EventLink
	completedJobs int
	totalJobs     int
}

func NewScraper() *Scraper {
	return &Scraper{
		collector: createCollector(),
		store:     models.NewEventStore(),
	}
}

func createCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(config.DOMAIN),
		colly.MaxDepth(config.MAX_DEPTH),
		colly.Async(config.ASYNC),
		colly.AllowURLRevisit(),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  config.GLOB,
		Parallelism: config.WORKERS,
		Delay:       config.REQUEST_DELAY * time.Second,
		RandomDelay: config.REQUEST_RANDOM_DELAY * time.Second,
	})

	return c
}

func (s *Scraper) Store() *models.EventStore {
	return s.store
}

func (s *Scraper) ScrapeEventList(startURL string) error {
	s.registerListHandlers()

	fmt.Printf("Scraping event list from: %s\n", startURL)
	if err := s.collector.Visit(startURL); err != nil {
		return fmt.Errorf("failed to visit event list: %w", err)
	}

	s.collector.Wait()

	fmt.Printf("Found and stored %d event URLs\n", len(s.eventURLs))
	return nil
}

func (s *Scraper) ScrapeEventDetails() error {
	if len(s.eventURLs) == 0 {
		return fmt.Errorf("no event URLs to scrape")
	}

	s.clearHandlers()
	s.registerDetailHandlers()

	s.totalJobs = len(s.eventURLs)
	s.completedJobs = 0

	startTime := time.Now()

	for i, link := range s.eventURLs {
		ctx := colly.NewContext()
		ctx.Put("eventName", link.Name)
		ctx.Put("jobIndex", fmt.Sprintf("%d", i+1))

		if err := s.collector.Request("GET", link.URL, nil, ctx, nil); err != nil {
			log.Printf("Failed to queue request for %s: %v", link.URL, err)
		}
	}

	fmt.Printf("All %d requests are queued, processing...\n", s.totalJobs)

	s.collector.Wait()

	if s.completedJobs == s.totalJobs {
		fmt.Println()
	}

	duration := time.Since(startTime)
	fmt.Printf("\rCompleted %d/%d jobs in %v\033[K\n",
		s.completedJobs, s.totalJobs, duration.Truncate(time.Second))

	return nil
}

func (s *Scraper) registerListHandlers() {
	s.collector.OnHTML(SelectorEventsTable, s.handleEventTable)
	s.collector.OnError(func(r *colly.Response, err error) {
		fmt.Printf("List error for %s: %v\n", r.Request.URL, err)
	})
}

func (s *Scraper) registerDetailHandlers() {
	s.collector.OnHTML(SelectorEventHeader, s.handleEventDetails)

	s.collector.OnError(func(r *colly.Response, err error) {
		eventName := r.Ctx.Get("eventName")
		fmt.Printf("\rFailed %s - %v\n", eventName, err)
		s.incrementAndShowProgress("", 0)
	})
}

func (s *Scraper) handleEventTable(el *colly.HTMLElement) {
	el.ForEach("tr", func(_ int, tr *colly.HTMLElement) {
		event, err := parseEventRow(tr)
		if err != nil {
			log.Println("Skipping invalid row: ", err)
			return
		}

		s.store.AddEvent(event)

		if event.URL != "" {
			s.eventURLs = append(s.eventURLs, EventLink{
				Name: event.EventName,
				URL:  event.URL,
			})
		}
	})
}

func (s *Scraper) handleEventDetails(el *colly.HTMLElement) {
	eventName := el.Request.Ctx.Get("eventName")

	eventDetails, found, err := parseEventDetails(el)
	if err != nil {
		fmt.Printf("Error parsing %s: %v\n", eventName, err)
		s.incrementAndShowProgress("", 0)
		return
	}

	if !found {
		return
	}

	if len(eventDetails.Outcomes) == 0 {
		s.incrementAndShowProgress("", 0)
		return
	}

	s.store.AddEventDetail(eventName, eventDetails)
	s.incrementAndShowProgress(eventName, len(eventDetails.Outcomes))
}

func (s *Scraper) incrementAndShowProgress(eventName string, outcomes int) {
	s.jobMutex.Lock()
	defer s.jobMutex.Unlock()

	s.completedJobs++

	var message string
	if eventName != "" && outcomes > 0 {
		message = fmt.Sprintf("[%d/%d] Latest: %s (%d outcomes)",
			s.completedJobs, s.totalJobs, eventName, outcomes)
	} else {
		message = fmt.Sprintf("[%d/%d] Processing...",
			s.completedJobs, s.totalJobs)
	}

	fmt.Printf("\r%s\033[K", message)
}

func (s *Scraper) clearHandlers() {
	s.collector = createCollector()
}
