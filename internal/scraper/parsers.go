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

	switch h3text {
	case EventChoiceHeader:
		// Choices & Outcomes
		return parseChoiceEvent(el)
	case EventNonChoiceHeader:
		// No-Choice Event
		return parseNoChoiceEvent(el)
	case EventAcupuncturistEventHeader:
		// Acupuncturist Event
		return parseAcupuncturistEvent(el)
	case EventNonChoiceAndOutcomesHeader:
		// No Choices and Outcomes
		return parseNoChoiceNoOutcomeEvent()
	default:
		// Couldn't locate header
		return models.EventDetails{}, false, nil
	}
}

func parseChoiceEvent(el *colly.HTMLElement) (models.EventDetails, bool, error) {
	table := el.DOM.NextFiltered(SelectorEventTable)
	if table.Length() == 0 {
		return models.EventDetails{}, true,
			fmt.Errorf("no table found after %q", EventChoiceHeader)
	}

	var details models.EventDetails
	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		if i == headerRowIndex {
			return // skip headaer
		}

		label := cleanHTMLTextSingleLine(row.Find("td:nth-child(1)"))
		reward := cleanHTMLTextWithBreaks(row.Find("td:nth-child(2)"))

		details.Outcomes = append(details.Outcomes, models.Choice{
			Label:  label,
			Reward: reward,
		})
	})

	return details, true, nil
}

func parseNoChoiceEvent(el *colly.HTMLElement) (models.EventDetails, bool, error) {
	table := el.DOM.NextFiltered(SelectorEventTable)
	if table.Length() == 0 {
		return models.EventDetails{}, true,
			fmt.Errorf("no table found after %q", EventNonChoiceHeader)
	}

	reward := cleanHTMLTextWithBreaks(table.Find("tbody tr").Eq(1).Find("td"))

	return models.EventDetails{
		Outcomes: []models.Choice{{
			Label:  "No-Choices",
			Reward: reward,
		}},
	}, true, nil
}

func parseAcupuncturistEvent(el *colly.HTMLElement) (models.EventDetails, bool, error) {
	table := el.DOM.NextAllFiltered(SelectorEventTable)
	if table.Length() == 0 {
		return models.EventDetails{}, true,
			fmt.Errorf("no table found after %q", EventAcupuncturistEventHeader)
	}

	var details models.EventDetails
	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		if i == headerRowIndex {
			return
		}

		label := cleanHTMLTextSingleLine(row.Find("th"))
		reward := cleanHTMLTextWithBreaks(row.Find("td"))

		details.Outcomes = append(details.Outcomes, models.Choice{
			Label:  label,
			Reward: reward,
		})
	})

	return details, true, nil
}

func parseNoChoiceNoOutcomeEvent() (models.EventDetails, bool, error) {
	return models.EventDetails{
		Outcomes: []models.Choice{{
			Label:  "No-Choices",
			Reward: "This scenario event contains no meaningful choices and does not affect your trainee's stats, mood, condition, or provide any skill hints.",
		}},
	}, true, nil
}

func cleanHTMLTextWithBreaks(s *goquery.Selection) string {
	clone := s.Clone()
	clone.Find("br").ReplaceWithHtml("\n")
	clone.Find("hr").Remove()

	text := clone.Text()

	lines := strings.Split(text, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	return strings.Join(cleanLines, "\n")
}

func cleanHTMLTextSingleLine(s *goquery.Selection) string {
	clone := s.Clone()
	clone.Find("hr").Remove()
	clone.Find("br").Remove()

	text := clone.Text()
	text = strings.Join(strings.Fields(text), " ")

	return text
}
