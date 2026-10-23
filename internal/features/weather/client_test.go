package weather

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_testutils "github.com/moond0wner/morning-digest/internal/core/testutils"
)

func TestClient_GetCoordinates(t *testing.T) {
	t.Run("обработка данных", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`[
			{"name": "Саратов", "lat": 51.54, "lon": 45.99, "country": "RU"}
		]`))
		}))
		defer server.Close()

		c := newClient("fake_api_key", server.URL, server.URL)

		geo, err := c.getCoordinates(context.Background(), "Саратов")
		if err != nil {
			t.Fatalf("getCoordinates вернул ошибку: %v", err)
		}

		if len(geo) != 1 {
			t.Fatalf("ожидали 1 результат, получили %d", len(geo))
		}
		core_testutils.AssertEqual(t, geo[0].Name, "Саратов")
		core_testutils.AssertEqual(t, geo[0].Lat, 51.54)
	})

	t.Run("пустой массив данных", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`[]`))
		}))
		defer server.Close()

		c := newClient("fake_api_key", server.URL, server.URL)

		_, err := c.getCoordinates(context.Background(), "Саратов")
		if err == nil {
			t.Fatal("ожидали ошибку, а получили nil")

		}
		targetErr := &core_errors.Error{Type: core_errors.ErrNotFound}
		if !errors.Is(err, targetErr) {
			t.Fatalf("ожидали ошибку ErrNotFound, получили: %v", err)
		}
	})

	t.Run("ошибка сервера http.StatusInternalServerError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`[{code: 500}]`))
		}))
		defer server.Close()
		c := newClient("fake_api_key", server.URL, server.URL)
		_, err := c.getCoordinates(context.Background(), "Саратов")
		if err == nil {
			t.Fatal("ожидали ошибку, а получили nil")

		}
		targetErr := &core_errors.Error{Type: core_errors.ErrExternalAPI}
		if !errors.Is(err, targetErr) {
			t.Fatalf("ожидали ошибку ErrExternalAPI, получили: %v", err)
		}
	})

	t.Run("битые координаты", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`[
			{"name": "билиберда", "lat": "abcdef", "lon": "abcdef", "country": "BEBEBE"}
		]`))
		}))
		defer server.Close()
		c := newClient("fake_api_key", server.URL, server.URL)
		_, err := c.getCoordinates(context.Background(), "Саратов")
		if err == nil {
			t.Fatal("ожидали ошибку, а получили nil")
		}

		targetErr := &core_errors.Error{Type: core_errors.ErrInternal}
		if !errors.Is(err, targetErr) {
			t.Fatalf("ожидали ошибку ErrInternal, получили: %v", err)
		}
	})
}
