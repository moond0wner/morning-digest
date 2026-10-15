package habr

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
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
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"create parse request",
		)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrExternalAPI,
			err,
			"do parse request",
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, core_errors.Wrap(
			core_errors.CheckStatusCode(response.StatusCode),
			fmt.Errorf("status %s", response.Status),
			"response status code",
		)
	}

	var rss RSS
	if err := xml.NewDecoder(response.Body).Decode(&rss); err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"xml decode",
		)
	}

	if len(rss.Channel.Items) < limit {
		limit = len(rss.Channel.Items)
	}
	top := rss.Channel.Items[:limit]

	return top, nil
}
