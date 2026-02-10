package duckduckgo

// External DataSource Adapter for DuckDuckGo HTML search
import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"strings"
	"time"

	"locus/models"

	"github.com/PuerkitoBio/goquery"
)

const defaultQuestionCount = 5

type DataSourceDuckDuckGo struct {
	Client     *http.Client
	BaseURL    string
	UserAgent  string
	SiteFilter string
	Debug      bool // Print lightweight fetch diagnostics when true
}

func New() *DataSourceDuckDuckGo {
	return &DataSourceDuckDuckGo{
		Client: &http.Client{
			Timeout: 8 * time.Second,
		},
		BaseURL:    "https://duckduckgo.com/html/",
		UserAgent:  "locus/duckduckgo-datasource",
		SiteFilter: "",
	}
}

// Init implements models.DataSource. DuckDuckGo requires no heavy initialization
func (es *DataSourceDuckDuckGo) Init() error {
	if es.Client == nil {
		es.Client = &http.Client{Timeout: 8 * time.Second}
	}
	if es.BaseURL == "" {
		es.BaseURL = "https://duckduckgo.com/html/"
	}
	if es.UserAgent == "" {
		es.UserAgent = "locus/duckduckgo-datasource"
	}
	return nil
}

// CheckAvailability implements models.DataSource
// Performs a lightweight search request to verify connectivity and expected response structure
func (es *DataSourceDuckDuckGo) CheckAvailability() bool {
	if err := es.Init(); err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	searchURL := es.buildSearchURL("duckduckgo")
	resp, err := es.doRequest(ctx, searchURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// FetchTopics implements models.DataSource
func (es *DataSourceDuckDuckGo) FetchTopics(count int, input string) ([]models.DataSourceTopic, error) {
	query := strings.TrimSpace(input)
	if query == "" {
		return nil, errors.New("Missing Search Input for DuckDuckGo data source")
	}
	if count <= 0 {
		count = defaultQuestionCount
	}
	if err := es.Init(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	searchURL := es.buildSearchURL(query)
	if es.Debug {
		fmt.Printf("[duckduckgo] search url: %s\n", searchURL)
	}
	resp, err := es.doRequest(ctx, searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("duckduckgo request failed: status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	if es.Debug {
		pageTitle := strings.TrimSpace(doc.Find("title").First().Text())
		fmt.Printf("[duckduckgo] page title: %s\n", pageTitle)
	}

	results := make([]models.DataSourceTopic, 0, count)
	seen := map[string]struct{}{}

	// DuckDuckGo markup can vary, so keep the primary selector broad
	selector := "a.result__a, a.result__a.js-result-title-link, a.result__url"
	selection := doc.Find(selector)
	if es.Debug {
		fmt.Printf("[duckduckgo] selector matches: %d\n", selection.Length())
	}
	selection.EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if len(results) >= count {
			return false
		}

		title := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")
		resolved := es.normalizeResultURL(strings.TrimSpace(href))
		if title == "" || resolved == "" {
			return true
		}
		if _, ok := seen[resolved]; ok {
			return true
		}
		seen[resolved] = struct{}{}

		results = append(results, models.DataSourceTopic{
			Topic:   normalizeWhitespace(title),
			SourceURL:  resolved,
			TopicID: urlToID(resolved),
			Site:       "duckduckgo",
		})
		return true
	})

	// If standard anchors are missing, fall back to a site-filtered scan
	if len(results) == 0 {
		results = es.fallbackResultLinks(doc, count, seen)
		if es.Debug {
			fmt.Printf("[duckduckgo] fallback results: %d\n", len(results))
		}
	}

	if len(results) == 0 {
		return nil, nil
	}
	return results, nil
}

// FetchData implements models.DataSource. 
// DuckDuckGo does not provide a way to fetch detailed data for a topic, so this is a no-op.
func (es *DataSourceDuckDuckGo) FetchData(count int, topicID int64) ([]models.DataSourceData, error) {
	return []models.DataSourceData{}, nil
}

// buildSearchURL constructs the DuckDuckGo search URL with the given query and site filter if set.
func (es *DataSourceDuckDuckGo) buildSearchURL(query string) string {
	base := strings.TrimRight(es.BaseURL, "/")
	values := url.Values{}
	values.Set("q", es.buildQuery(query))
	return fmt.Sprintf("%s/?%s", base, values.Encode())
}

// buildQuery constructs the search query string, applying the site filter if configured.
func (es *DataSourceDuckDuckGo) buildQuery(query string) string {
	filter := strings.TrimSpace(es.SiteFilter)
	if filter == "" {
		return query
	}
	// Normalize to a `site:` modifier so callers can pass either form
	if strings.HasPrefix(filter, "site:") {
		return fmt.Sprintf("%s %s", filter, query)
	}
	return fmt.Sprintf("site:%s %s", filter, query)
}

// fallbackResultLinks performs a broad scan of all anchor tags in the document to find links matching the site filter.
func (es *DataSourceDuckDuckGo) fallbackResultLinks(doc *goquery.Document, count int, seen map[string]struct{}) []models.DataSourceTopic {
	targetHost := strings.TrimSpace(es.SiteFilter)
	if targetHost == "" {
		return nil
	}
	if strings.HasPrefix(targetHost, "site:") {
		targetHost = strings.TrimSpace(strings.TrimPrefix(targetHost, "site:"))
	}
	if targetHost == "" {
		return nil
	}

	results := make([]models.DataSourceTopic, 0, count)
	// Scan all anchors and keep only matches for the target host
	doc.Find("a[href]").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if len(results) >= count {
			return false
		}
		text := strings.TrimSpace(s.Text())
		href, _ := s.Attr("href")
		resolved := es.normalizeResultURL(strings.TrimSpace(href))
		if resolved == "" {
			return true
		}
		parsed, err := url.Parse(resolved)
		if err != nil || parsed.Host == "" {
			return true
		}
		if !strings.HasSuffix(parsed.Host, targetHost) {
			return true
		}
		if _, ok := seen[resolved]; ok {
			return true
		}
		seen[resolved] = struct{}{}

		title := text
		if title == "" {
			title = resolved
		}
		results = append(results, models.DataSourceTopic{
			Topic:   normalizeWhitespace(title),
			SourceURL:  resolved,
			TopicID: urlToID(resolved),
			Site:       "duckduckgo",
		})
		return true
	})
	return results
}

// normalizeResultURL processes a raw URL from DuckDuckGo search results, resolving relative URLs and filtering out ad links.
func (es *DataSourceDuckDuckGo) normalizeResultURL(raw string) string {
	if raw == "" {
		return ""
	}
	// Skip ad links early to avoid polluting results
	if strings.Contains(raw, "ad_domain") {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if !parsed.IsAbs() {
		base, err := url.Parse(es.BaseURL)
		if err == nil {
			parsed = base.ResolveReference(parsed)
		}
	}
	if strings.Contains(parsed.Host, "duckduckgo.com") && strings.HasPrefix(parsed.Path, "/l/") {
		if target := parsed.Query().Get("uddg"); target != "" {
			if decoded, err := url.QueryUnescape(target); err == nil {
				return decoded
			}
			return target
		}
	}
	// Drop links that were tagged as ads after redirect resolution
	if parsed.Query().Has("ad_domain") {
		return ""
	}
	return parsed.String()
}

// doRequest performs an HTTP GET request to the specified URL with appropriate headers and context.
func (es *DataSourceDuckDuckGo) doRequest(ctx context.Context, target string) (*http.Response, error) {
	client := es.Client
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html")
	if es.UserAgent != "" {
		req.Header.Set("User-Agent", es.UserAgent)
	}
	return client.Do(req)
}

// Helpers
func urlToID(raw string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(raw))
	return int64(h.Sum64())
}
func normalizeWhitespace(in string) string {
	fields := strings.Fields(in)
	return strings.Join(fields, " ")
}
