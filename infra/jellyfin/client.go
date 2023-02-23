package jellyfin

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ClientConfig struct {
	BaseURL string
	ApiKey  string
	UserID  string
	Headers map[string]string
}

type Client struct {
	client http.Client
	config ClientConfig
}

func (c *ClientConfig) sanitize() {
	c.BaseURL = strings.TrimRight(c.BaseURL, "/")

	c.Headers = make(map[string]string)
	// https://gist.github.com/nielsvanvelzen/ea047d9028f676185832e51ffaf12a6f
	c.Headers["Authorization"] = fmt.Sprintf("MediaBrowser Token=\"%s\"", c.ApiKey)
}

func NewClient(config ClientConfig) *Client {
	config.sanitize()

	client := http.Client{}

	return &Client{
		client: client,
		config: config,
	}
}

func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.config.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}
