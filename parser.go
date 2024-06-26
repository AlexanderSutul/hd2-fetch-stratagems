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

type HelldiverParser struct {
	baseUrl string
}

var _ Parser = (*HelldiverParser)(nil)

func NewParser(baseUrl string) *HelldiverParser {
	return &HelldiverParser{baseUrl: baseUrl}
}

func (h *HelldiverParser) Parse() ([]*Stratagem, error) {
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

	stratagems := h.getStratagemsFromDocument(doc)

	return stratagems, nil
}

func (h *HelldiverParser) makeRequest() (*http.Response, error) {
	req, err := http.NewRequest("GET", h.baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	client := http.Client{Timeout: 5 * time.Second}
	return client.Do(req)
}

func (h *HelldiverParser) getStratagemsFromDocument(doc *goquery.Document) []*Stratagem {
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

func (h *HelldiverParser) getIconUrl(s *goquery.Selection) string {
	val, _ := s.Find("a > img").Attr("data-src")
	return val
}

func (h *HelldiverParser) getName(s *goquery.Selection) string {
	return s.Find("a").Text()
}

func (h *HelldiverParser) getInputCodeFromSelector(s *goquery.Selection) InputCode {
	alt, _ := s.Attr("alt")
	spl := strings.Split(alt, " ")
	altPostfix := spl[len(spl)-1]
	l := strings.ToLower(altPostfix)
	m := map[string]InputCode{
		"u": UP,
		"d": DOWN,
		"l": LEFT,
		"r": RIGHT,
	}
	return m[l]
}

func (h *HelldiverParser) getInputs(s *goquery.Selection) []InputCode {
	var inputs []InputCode

	s.Find("img").Each(func(i int, s *goquery.Selection) {
		inputs = append(inputs, h.getInputCodeFromSelector(s))
	})

	return inputs
}

func (h *HelldiverParser) getCooldown(s *goquery.Selection) string {
	return strings.TrimRight(s.Text(), "\n")
}

func (h *HelldiverParser) getUses(s *goquery.Selection) string {
	return strings.TrimRight(s.Text(), "\n")
}

func (h *HelldiverParser) getActivation(s *goquery.Selection) string {
	return strings.TrimRight(s.Text(), "\n")
}
