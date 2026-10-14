package habr

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type client struct {
	http *http.Client
}

func newClient() *client {
	return &client{
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *client) parse(ctx context.Context, url string, limit int) ([]Item, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create parse request: %w", err)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do parse request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %d", response.StatusCode)
	}

	var rss RSS
	if err := xml.NewDecoder(response.Body).Decode(&rss); err != nil {
		return nil, fmt.Errorf("xml decode: %w", err)
	}

	if len(rss.Channel.Items) < limit {
		limit = len(rss.Channel.Items)
	}
	top := rss.Channel.Items[:limit]

	return top, nil
}
