package scraper

import (
	"github.com/arthwr/uma-scrap/internal/models"
	"github.com/gocolly/colly"
)

func parseEventRow(tr *colly.HTMLElement) (models.Event, error) {
	umaName := tr.ChildText("td:nth-child(1) a")
	if umaName == "" {
		umaName = tr.ChildText("td:nth-child(1) div")
	}
	eventTypeStr := tr.ChildText("td:nth-child(2)")
	eventName := tr.ChildText("td:nth-child(3) a")
	eventURL := tr.ChildAttr("td:nth-child(3) a", "href")

	t, err := models.EventTypeFromString(eventTypeStr)
	if err != nil {
		return models.Event{}, err
	}

	return models.Event{
		UmaName:   umaName,
		EventName: eventName,
		URL:       eventURL,
		Type:      t,
	}, nil
}
