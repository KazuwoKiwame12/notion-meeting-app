package function

import (
	"app/config"
	"app/domain/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.notion.com/v1/"
const apiVersion = "2021-05-13"

type NotionClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewNotionClient() *NotionClient {
	nc := &NotionClient{
		apiKey:     config.NotionToken(),
		httpClient: http.DefaultClient,
	}
	return nc
}

func (nc *NotionClient) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, baseURL+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", nc.apiKey))
	req.Header.Set("Notion-Version", apiVersion)

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
		return fmt.Errorf("notion: failed to create page: %v", res)
	}

	return nil
}
