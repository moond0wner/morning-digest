package habr

import (
	"context"
	"fmt"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
)

type service struct {
	client *client
	log    *core_logger.Logger
	retry  *core_retry.Retrier
}

func newService(logger *core_logger.Logger, retrier *core_retry.Retrier) *service {
	return &service{
		client: newClient(),
		log:    logger,
		retry:  retrier,
	}
}

const (
	baseRSSUrl = "https://habr.com/ru/rss/feed/cf9c3edefc5fc024bbc10312c7f02eab/?fl=ru"
	hubRSSUrl  = "https://habr.com/ru/rss/hubs/%s/articles/?fl=ru"
)

func (s *service) baseParse(ctx context.Context, limit int) ([]Article, error) {
	articles, err := s.parseFromUrl(ctx, baseRSSUrl, limit)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *service) hubParse(ctx context.Context, hub string, limit int) ([]Article, error) {
	hubUrl := fmt.Sprintf(hubRSSUrl, hub)

	articles, err := s.parseFromUrl(ctx, hubUrl, limit)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (s *service) parseFromUrl(ctx context.Context, url string, limit int) ([]Article, error) {
	var articles []Article
	if err := s.retry.Do(ctx, func(ctx context.Context) error {
		items, err := s.client.parse(ctx, url, limit)
		if err != nil {
			return err
		}
		articles, err = convertToArticle(items)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("retry error: %w", err)
	}

	return articles, nil
}

func convertToArticle(items []Item) ([]Article, error) {
	articles := make([]Article, 0, len(items))
	for _, item := range items {
		article, err := itemToArticle(item)
		if err != nil {
			return nil, core_errors.Wrap(
				core_errors.ErrInternal,
				err,
				"failed turn item into article",
			)
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func timeParse(value string) (time.Time, error) {
	layout := time.RFC1123

	result, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"time parse",
		)
	}
	return result, nil
}

func itemToArticle(item Item) (Article, error) {
	t, err := timeParse(item.PubDate)
	if err != nil {
		return Article{}, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"item to article",
		)
	}
	return Article{
		Title: item.Title,
		Link:  item.Link,
		Date:  t,
	}, nil
}
