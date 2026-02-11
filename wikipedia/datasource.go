package wikipedia

// Data Source Adapter for Wikipedia API
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	datasource "github.com/locus-search/datasource-sdk"
)

type DataSourceWikipedia struct {
	Client    *http.Client
	BaseURL   string
	UserAgent string
}

func New() *DataSourceWikipedia {
	return &DataSourceWikipedia{
		Client: &http.Client{
			Timeout: 8 * time.Second,
		},
		BaseURL:   "https://en.wikipedia.org/w/api.php",
		UserAgent: "locus/ask",
	}
}

// Init implements models.DataSource
// Wikipedia requires no initialization
func (es *DataSourceWikipedia) Init() error {
	return nil
}

// CheckAvailability implements models.DataSource
func (es *DataSourceWikipedia) CheckAvailability() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	params := url.Values{}
	params.Set("action", "query")
	params.Set("meta", "siteinfo")
	params.Set("format", "json")

	_, err := es.doJSON(ctx, params, &struct{}{})
	return err == nil
}

// FetchTopics implements models.DataSource
// Fetch Wikipedia search results for the query string. Each result is a topic with title and page ID.
func (es *DataSourceWikipedia) FetchTopics(count int, input string) ([]datasource.DataSourceTopic, error) {
	query := strings.TrimSpace(input)
	if query == "" {
		return nil, errors.New("Missing search input for Wikipedia DataSource")
	}
	if count <= 0 {
		count = 5
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	params := url.Values{}
	params.Set("action", "query")
	params.Set("list", "search")
	params.Set("srsearch", query)
	params.Set("srlimit", fmt.Sprintf("%d", count))
	params.Set("format", "json")

	var response struct {
		Query struct {
			Search []struct {
				Title  string `json:"title"`
				PageID int64  `json:"pageid"`
			} `json:"search"`
		} `json:"query"`
		Error *struct {
			Info string `json:"info"`
		} `json:"error"`
	}

	_, err := es.doJSON(ctx, params, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("wikipedia error: %s", response.Error.Info)
	}

	results := make([]datasource.DataSourceTopic, 0, len(response.Query.Search))
	for _, item := range response.Query.Search {
		results = append(results, datasource.DataSourceTopic{
			Topic:   item.Title,
			SourceURL:  fmt.Sprintf("https://en.wikipedia.org/?curid=%d", item.PageID),
			TopicID: item.PageID,
		})
	}
	return results, nil
}

// FetchData implements models.DataSource
// Fetch the extract (intro paragraph) for the given Wikipedia page ID
// Returns a single DataSourceData item with the extract text and source URL
func (es *DataSourceWikipedia) FetchData(count int, topicID int64) ([]datasource.DataSourceData, error) {
	if topicID <= 0 {
		return nil, errors.New("topicID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	params := url.Values{}
	params.Set("action", "query")
	params.Set("pageids", fmt.Sprintf("%d", topicID))
	params.Set("prop", "extracts")
	params.Set("exintro", "1")
	params.Set("explaintext", "1")
	params.Set("format", "json")

	var response struct {
		Query struct {
			Pages map[string]struct {
				PageID  int64  `json:"pageid"`
				Title   string `json:"title"`
				Extract string `json:"extract"`
			} `json:"pages"`
		} `json:"query"`
		Error *struct {
			Info string `json:"info"`
		} `json:"error"`
	}

	_, err := es.doJSON(ctx, params, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, fmt.Errorf("wikipedia error: %s", response.Error.Info)
	}

	for _, page := range response.Query.Pages {
		dataText := strings.TrimSpace(page.Extract)
		if dataText == "" {
			return []datasource.DataSourceData{}, nil
		}
		data := datasource.DataSourceData{
			DataText: dataText,
			SourceURL:  fmt.Sprintf("https://en.wikipedia.org/?curid=%d", page.PageID),
			AnswerID:   page.PageID,
		}
		return []datasource.DataSourceData{data}, nil
	}

	return []datasource.DataSourceData{}, nil
}

// doJSON performs an HTTP GET request to the Wikipedia API with the specified parameters and decodes the JSON response into the target structure
func (es *DataSourceWikipedia) doJSON(ctx context.Context, params url.Values, target interface{}) (int, error) {
	client := es.Client
	if client == nil {
		client = &http.Client{Timeout: 8 * time.Second}
	}
	endpoint := strings.TrimRight(es.BaseURL, "/")
	uri := endpoint
	if encoded := params.Encode(); encoded != "" {
		uri = uri + "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return 0, err
	}
	if es.UserAgent != "" {
		req.Header.Set("User-Agent", es.UserAgent)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, fmt.Errorf("wikipedia request failed: %s", strings.TrimSpace(string(body)))
	}

	if target == nil {
		return resp.StatusCode, nil
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, nil
}
