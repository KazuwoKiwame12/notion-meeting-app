package function

import (
	"app/domain/model"
	"app/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type NotionClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewNotionClient(apiKey string) *NotionClient {
	nc := &NotionClient{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}
	return nc
}

func (nc *NotionClient) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, config.NOTION_API_URL()+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", nc.apiKey))
	req.Header.Set("Notion-Version", config.NOTION_API_VERSION())

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (nc *NotionClient) CreatePage(params model.Template) error {
	body := &bytes.Buffer{}

	err := json.NewEncoder(body).Encode(params)
	if err != nil {
		return fmt.Errorf("notion: failed to encode body params to JSON: %w", err)
	}

	req, err := nc.newRequest(http.MethodPost, "/pages", body)
	if err != nil {
		return fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := nc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		result, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("notion: failed to create page: %s", result)
	}

	return nil
}
