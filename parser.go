package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
)

const (
	SectionClass = "mw-headline"
)

type HelldiversParser struct {
	baseUrl string
}

func NewParser(baseUrl string) *HelldiversParser {
	return &HelldiversParser{baseUrl: baseUrl}
}

func (h *HelldiversParser) Parse() (GroupedStratagems, error) {
	res, err := h.makeRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s: %w", h.baseUrl, err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get %s: %s", h.baseUrl, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	stratagems := h.getStratagems(doc)

	return stratagems, nil
}

func (h *HelldiversParser) makeRequest() (*http.Response, error) {
	req, err := http.NewRequest("GET", h.baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	client := http.Client{Timeout: 5 * time.Second}
	return client.Do(req)
}

func (h *HelldiversParser) getStratagems(doc *goquery.Document) GroupedStratagems {
	var stratagems []*Stratagem
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		sectionName := s.Find("span").First().Text()

		s.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
			stratagem := &Stratagem{Group: sectionName}

			s.Find("td").Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0: // icon url
					stratagem.IconUrl = h.getIconUrl(s)
				case 1: // name
					stratagem.Name = h.getName(s)
				case 2: // inputs
					stratagem.InputCode = h.getInputs(s)
				case 3: // cooldown
					stratagem.Cooldown = h.getCooldown(s)
				case 4: // uses
					stratagem.Uses = h.getUses(s)
				case 5: // activation
					stratagem.Activation = h.getActivation(s)
				}

				stratagems = append(stratagems, stratagem)
			})
		})
	})
	return h.groupStratagemByGroup(stratagems)
}

func (h *HelldiversParser) groupStratagemByGroup(stratagems []*Stratagem) map[string][]*Stratagem {
	grouped := make(map[string][]*Stratagem)
	for _, stratagem := range stratagems {
		grouped[stratagem.Group] = append(grouped[stratagem.Group], stratagem)
	}
	return grouped
}

func (h *HelldiversParser) getIconUrl(s *goquery.Selection) string {
	val, _ := s.Find("a").Attr("href")
	return val
}

func (h *HelldiversParser) getName(s *goquery.Selection) string {
	return s.Find("a").Text()
}

func (h *HelldiversParser) getInputs(s *goquery.Selection) []InputCode {
	var inputs []InputCode

	s.Find("img").Each(func(i int, s *goquery.Selection) {
		alt, _ := s.Attr("alt")
		spl := strings.Split(alt, " ")
		altPostfix := spl[len(spl)-1]
		l := strings.ToLower(altPostfix)
		switch l {
		case "u":
			inputs = append(inputs, UP)
		case "d":
			inputs = append(inputs, DOWN)
		case "l":
			inputs = append(inputs, LEFT)
		case "r":
			inputs = append(inputs, RIGHT)
		}
	})

	return inputs
}

func (h *HelldiversParser) getCooldown(s *goquery.Selection) string {
	return strings.TrimRight(s.Text(), "\n")
}

func (h *HelldiversParser) getUses(s *goquery.Selection) string {
	return strings.TrimRight(s.Text(), "\n")
}

func (h *HelldiversParser) getActivation(s *goquery.Selection) string {
	return strings.TrimRight(s.Text(), "\n")
}