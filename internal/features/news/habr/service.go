package habr

import (
	"context"
	"fmt"
	"time"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
)

type Service struct {
	client *client
	logger *core_logger.Logger
}

func NewService() *Service {
	return &Service{
		client: newClient(),
	}
}

const (
	baseRSSUrl = "https://habr.com/ru/rss/feed/cf9c3edefc5fc024bbc10312c7f02eab/?fl=ru"
	hubRSSUrl  = "https://habr.com/ru/rss/hubs/%s/articles/?fl=ru"
)

func (s *Service) BaseParse(ctx context.Context, limit int) ([]Article, error) {
	items, err := s.client.parse(ctx, baseRSSUrl, limit)
	if err != nil {
		return nil, fmt.Errorf("base parse: %w", err)
	}

	articles := make([]Article, 0, len(items))
	for _, item := range items {
		article, err := itemToArticle(item)
		if err != nil {
			return nil, fmt.Errorf("failed turn item into article: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func (s *Service) HubParse(ctx context.Context, hub string, limit int) ([]Article, error) {
	hubUrl := fmt.Sprintf(hubRSSUrl, hub)
	items, err := s.client.parse(ctx, hubUrl, limit)
	if err != nil {
		return nil, fmt.Errorf("base parse: %w", err)
	}

	articles := make([]Article, 0, len(items))
	for _, item := range items {
		article, err := itemToArticle(item)
		if err != nil {
			return nil, fmt.Errorf("failed turn item into article: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func timeParse(value string) (time.Time, error) {
	layout := time.RFC1123

	result, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("time parse: %w", err)
	}
	return result, nil
}

func itemToArticle(item Item) (Article, error) {
	time, err := timeParse(item.PubDate)
	if err != nil {
		return Article{}, fmt.Errorf("item to article: %w", err)
	}
	return Article{
		Title: item.Title,
		Link:  item.Link,
		Date:  time,
	}, nil
}
