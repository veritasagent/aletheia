package aletheia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type EvidencePack struct {
	Sources       []string
	Headlines     []string
	WikiSummary   string
	InstantAnswer string
	Raw           string
}

func FetchEvidence(claim string) EvidencePack {
	pack := EvidencePack{}
	newsAPIKey := os.Getenv("NEWS_API_KEY")

	// STEP 1 — NewsAPI
	if newsAPIKey != "" {
		q := claim
		if len(q) > 60 {
			q = q[:60]
		}
		newsURL := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&language=en&sortBy=publishedAt&pageSize=3&apiKey=%s",
			url.QueryEscape(q), newsAPIKey)

		resp, err := http.Get(newsURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var data struct {
				Status   string `json:"status"`
				Articles []struct {
					Title       string `json:"title"`
					Description string `json:"description"`
					Source      struct {
						Name string `json:"name"`
					} `json:"source"`
				} `json:"articles"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && data.Status == "ok" {
				for _, art := range data.Articles {
					pack.Headlines = append(pack.Headlines, fmt.Sprintf("%s (%s)", art.Title, art.Source.Name))
					pack.Sources = append(pack.Sources, art.Source.Name)
				}
			}
		}
	}

	// STEP 2 — Wikipedia
	words := strings.Fields(claim)
	stopWords := map[string]bool{
		"is": true, "the": true, "a": true, "an": true, "has": true,
		"have": true, "been": true, "all": true, "from": true,
	}
	var sigWord string
	count := 0
	for _, w := range words {
		lw := strings.ToLower(w)
		// Strip punctuation
		lw = strings.Trim(lw, ".,!?;:\"'")
		if !stopWords[lw] && lw != "" {
			sigWord = lw
			count++
			if count >= 1 { // We only need the first significant word per spec
				break
			}
		}
	}

	if sigWord != "" {
		wikiURL := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", url.PathEscape(sigWord))
		resp, err := http.Get(wikiURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var data struct {
				Extract string `json:"extract"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
				pack.WikiSummary = data.Extract
				if len(pack.WikiSummary) > 200 {
					pack.WikiSummary = pack.WikiSummary[:200]
				}
				pack.Sources = append(pack.Sources, "Wikipedia")
			}
		}
	}

	// STEP 3 — DuckDuckGo
	ddgURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_html=1&skip_disambig=1", url.QueryEscape(claim))
	resp, err := http.Get(ddgURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var data struct {
			AbstractText string `json:"AbstractText"`
			Answer       string `json:"Answer"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			if data.AbstractText != "" {
				pack.InstantAnswer = data.AbstractText
				if len(pack.InstantAnswer) > 150 {
					pack.InstantAnswer = pack.InstantAnswer[:150]
				}
			} else if data.Answer != "" {
				pack.InstantAnswer = data.Answer
			}
			if pack.InstantAnswer != "" {
				pack.Sources = append(pack.Sources, "DuckDuckGo")
			}
		}
	}

	// STEP 4 — Build Evidence Pack Raw
	if len(pack.Sources) == 0 {
		pack.Raw = "No corroborating sources found for this claim."
		return pack
	}

	var sb strings.Builder
	sb.WriteString("LIVE EVIDENCE RETRIEVED:\n\n")

	if len(pack.Headlines) > 0 {
		sb.WriteString("NewsAPI Headlines:\n")
		for _, h := range pack.Headlines {
			sb.WriteString(fmt.Sprintf("- %s\n", h))
		}
		sb.WriteString("\n")
	}

	if pack.WikiSummary != "" {
		sb.WriteString("Wikipedia Context:\n")
		sb.WriteString(pack.WikiSummary)
		sb.WriteString("\n\n")
	}

	if pack.InstantAnswer != "" {
		sb.WriteString("Instant Answer:\n")
		sb.WriteString(pack.InstantAnswer)
		sb.WriteString("\n")
	}

	pack.Raw = strings.TrimSpace(sb.String())
	return pack
}
