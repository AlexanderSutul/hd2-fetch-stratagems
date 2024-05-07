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

func (h *HelldiversParser) Parse() (Stratagems, error) {
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

func (h *HelldiversParser) getStratagems(doc *goquery.Document) Stratagems {
	var stratagems []*Stratagem
	doc.Find("table").Each(func(i int, tableSelect *goquery.Selection) {
		section := tableSelect.Find("span").First().Text()
		tableSelect.Find("tbody tr").Each(func(i int, tableRowSelect *goquery.Selection) {
			// jump over a redundant rows in the table
			if i < 2 {
				return
			}
			stratagem := &Stratagem{Group: section}
			tableRowSelect.Find("td").Each(func(i int, tableDataSelect *goquery.Selection) {
				switch i {
				case 0: // icon url
					stratagem.IconUrl = h.getIconUrl(tableDataSelect)
				case 1: // name
					stratagem.Name = h.getName(tableDataSelect)
				case 2: // inputs
					stratagem.InputCode = h.getInputs(tableDataSelect)
				case 3: // cooldown
					stratagem.Cooldown = h.getCooldown(tableDataSelect)
				case 4: // uses
					stratagem.Uses = h.getUses(tableDataSelect)
				case 5: // activation
					stratagem.Activation = h.getActivation(tableDataSelect)
				}
			})

			stratagems = append(stratagems, stratagem)
		})
	})
	return stratagems
}

func (h *HelldiversParser) getIconUrl(s *goquery.Selection) string {
	val, _ := s.Find("a > img").Attr("data-src")
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
