package finnhub

import (
	"net/http"
	"net/url"
)

const (
	finnhubBaseURL     = "https://finnhub.io/api/v1"
	finnhubTokenHeader = "X-Finnhub-Token"
)

type Client struct {
	apiKey     string
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(key string) *Client {
	url, err := url.Parse(finnhubBaseURL)
	if err != nil {
		url = nil
	}
	return &Client{
		apiKey:     key,
		baseURL:    url,
		httpClient: &http.Client{},
	}
}

func (c *Client) Get(query string) (*http.Response, error) {
	request, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set(finnhubTokenHeader, c.apiKey)

	return c.httpClient.Do(request)
}
