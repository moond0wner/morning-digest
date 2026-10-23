package habr

import (
	"testing"
	"time"
)

func TestTimeParse(t *testing.T) {
	got, err := timeParse("Mon, 12 Oct 2026 12:00:00 GMT")
	if err != nil {
		t.Fatalf("timeParse: %v", err)
	}
	if got.Year() != 2026 || got.Month() != time.October || got.Day() != 12 {
		t.Errorf("неверная дата: %v", got)
	}

	_, err = timeParse("это не дата")
	if err == nil {
		t.Error("ожидали ошибку на невалидной дате")
	}
}

func TestConvertToArticle(t *testing.T) {
	items := []Item{
		{Title: "Статья 1", Link: "https://example.com/1", PubDate: "Mon, 12 Oct 2026 12:00:00 GMT"},
	}

	got, err := convertToArticle(items)
	if err != nil {
		t.Fatalf("convertToArticle: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("ожидали 1 статью, получили %d", len(got))
	}
	if got[0].Title != "Статья 1" {
		t.Errorf("ожидали 'Статья 1', получили %s", got[0].Title)
	}
}
