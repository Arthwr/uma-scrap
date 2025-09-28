package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/arthwr/uma-scrap/internal/models"
	"github.com/gocolly/colly"
)

const headerRowIndex = 0

func parseEventRow(el *colly.HTMLElement) (models.Event, error) {
	umaName := el.ChildText("td:nth-child(1) a")
	if umaName == "" {
		umaName = el.ChildText("td:nth-child(1) div")
	}
	eventTypeStr := el.ChildText("td:nth-child(2)")
	eventName := el.ChildText("td:nth-child(3) a")
	eventURL := el.ChildAttr("td:nth-child(3) a", "href")

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

func parseEventDetails(el *colly.HTMLElement) (models.EventDetails, bool, error) {
	h3text := strings.TrimSpace(el.Text)

	if h3text != EventHeaderText {
		// not the section we are looking for
		return models.EventDetails{}, false, nil
	}

	table := el.DOM.NextFiltered(SelectorEventTable)
	if table.Length() == 0 {
		return models.EventDetails{}, true, fmt.Errorf("no table found after %q", h3text)
	}

	var details models.EventDetails

	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		if i == headerRowIndex {
			return // skip header
		}

		label := strings.TrimSpace(row.Find("td:nth-child(1)").Text())
		reward := extractTextWithBreaks(row.Find("td:nth-child(2)"))

		details.Outcomes = append(details.Outcomes, models.Choice{
			Label:  label,
			Reward: reward,
		})
	})

	return details, true, nil
}

func extractTextWithBreaks(sel *goquery.Selection) string {
	var sb strings.Builder

	sel.Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "br" {
			sb.WriteString("\n")
		} else {
			sb.WriteString(s.Text())
		}
	})

	return strings.TrimSpace(sb.String())
}
