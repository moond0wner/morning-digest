package habr

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_testutils "github.com/moond0wner/morning-digest/internal/core/testutils"
)

func TestClient_Parse(t *testing.T) {
	t.Run("обработка данных", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`
			<rss>
			<channel>
				<item>
					<title>Тест</title>
					<link>https://example.com</link>
					<pubDate>Mon, 12 Oct 2026 12:00:00 GMT</pubDate>
				</item>
			</channel>
			</rss>
			`))
		}))
		defer server.Close()
		c := newClient()
		items, err := c.parse(context.Background(), server.URL, 1)
		if err != nil {
			t.Fatalf("parse вернул ошибку: %v", err)
		}

		item := items[0]
		core_testutils.AssertEqual(t, item.Title, "Тест")
		core_testutils.AssertEqual(t, item.Link, "https://example.com")
		core_testutils.AssertEqual(t, item.PubDate, "Mon, 12 Oct 2026 12:00:00 GMT")
	})
	t.Run("ошибка сервера http.StatusInternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`[{code: 500}]`))
		}))
		defer server.Close()
		c := newClient()
		_, err := c.parse(context.Background(), server.URL, 1)
		if err == nil {
			t.Fatal("ожидали ошибку, а получили nil")

		}
		targetErr := &core_errors.Error{Type: core_errors.ErrExternalAPI}
		if !errors.Is(err, targetErr) {
			t.Fatalf("ожидали ошибку ErrExternalAPI, получили: %v", err)
		}
	})

	t.Run("битый xml", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`
			"это вообще не xml"
			`))
		}))
		defer server.Close()
		c := newClient()
		_, err := c.parse(context.Background(), server.URL, 1)
		if err == nil {
			t.Fatal("ожидали ошибку, а получили nil")
		}

		targetErr := &core_errors.Error{Type: core_errors.ErrInternal}
		if !errors.Is(err, targetErr) {
			t.Fatalf("ожидали ошибку ErrInternal, получили: %v", err)
		}
	})
}
