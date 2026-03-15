package aletheia

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Ingester struct {
	watcher  *Watcher
	interval time.Duration
	running  bool
	seen     map[string]bool
	mu       sync.Mutex
}

func NewIngester(watcher *Watcher) *Ingester {
	return &Ingester{
		watcher:  watcher,
		interval: 120 * time.Second,
		seen:     make(map[string]bool),
	}
}

func (i *Ingester) Start() {
	i.mu.Lock()
	if i.running {
		i.mu.Unlock()
		return
	}
	i.running = true
	i.mu.Unlock()

	go func() {
		for {
			i.mu.Lock()
			if !i.running {
				i.mu.Unlock()
				return
			}
			i.mu.Unlock()

			claims := make([]struct {
				Title  string
				Source string
			}, 0)

			// Fetch from sources
			for _, t := range FetchGoogleNewsRSS() {
				claims = append(claims, struct{ Title, Source string }{t, "Google News"})
			}
			for _, t := range FetchRedditIndia() {
				claims = append(claims, struct{ Title, Source string }{t, "Reddit"})
			}

			for _, c := range claims {
				i.mu.Lock()
				if !i.running {
					i.mu.Unlock()
					return
				}

				key := strings.ToLower(strings.TrimSpace(c.Title))
				if i.seen[key] {
					i.mu.Unlock()
					continue
				}
				i.seen[key] = true
				i.mu.Unlock()

				event := MisinformationEvent{
					ID:        fmt.Sprintf("ING-%s", fmt.Sprintf("%d", time.Now().UnixNano())[:10]),
					Claim:     c.Title,
					Platform:  c.Source,
					Virality:  rand.Intn(490001) + 10000,
					Timestamp: time.Now(),
				}
				i.watcher.Emit(event)
				time.Sleep(3 * time.Second)
			}

			time.Sleep(i.interval)
		}
	}()
}

func (i *Ingester) Stop() {
	i.mu.Lock()
	i.running = false
	i.mu.Unlock()
}

func cleanTitle(t string) string {
	return strings.TrimSpace(html.UnescapeString(t))
}

func FetchGoogleNewsRSS() []string {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://news.google.com/rss/search?q=india+fake+news+misinformation&hl=en-IN&gl=IN&ceid=IN:en"

	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var rss struct {
		Channel struct {
			Items []struct {
				Title string `xml:"title"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil
	}

	titles := []string{}
	for _, item := range rss.Channel.Items {
		if strings.Contains(item.Title, "Google News") {
			continue
		}
		titles = append(titles, cleanTitle(item.Title))
		if len(titles) >= 3 {
			break
		}
	}
	return titles
}

func FetchRedditIndia() []string {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://www.reddit.com/r/india/new.json?limit=5", nil)
	req.Header.Set("User-Agent", "ALETHEIA/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var data struct {
		Data struct {
			Children []struct {
				Data struct {
					Title string `json:"title"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil
	}

	titles := []string{}
	for _, child := range data.Data.Children {
		titles = append(titles, cleanTitle(child.Data.Title))
		if len(titles) >= 3 {
			break
		}
	}
	return titles
}

func FetchPIBFactCheck() []string {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://pib.gov.in/RssMain.aspx?ModId=6&Lang=1&Regid=3"

	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var rss struct {
		Channel struct {
			Items []struct {
				Title string `xml:"title"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil
	}

	titles := []string{}
	for _, item := range rss.Channel.Items {
		titles = append(titles, cleanTitle(item.Title))
		if len(titles) >= 3 {
			break
		}
	}
	return titles
}
