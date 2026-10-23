package invest

import (
	"errors"
	"testing"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
	core_testutils "github.com/moond0wner/morning-digest/internal/core/testutils"
	"go.uber.org/zap"
)

type mockClient struct {
	portfolio *Portfolio
	err       error
}

func (c *mockClient) getPortfolio() (*Portfolio, error) {
	return c.portfolio, c.err
}

func (c *mockClient) close() error {
	return nil
}

func TestService_GetPortfolio(t *testing.T) {
	logger := &core_logger.Logger{
		Logger: zap.NewNop(),
	}
	retry := core_retry.NewRetrier(1, 1*time.Second, 1*time.Second, logger)
	t.Run("успешное получение", func(t *testing.T) {
		mock := &mockClient{
			portfolio: &Portfolio{
				TotalAmount:   10000,
				DailyYield:    500,
				ExpectedYield: 200,
			},
		}
		s := newService(mock, logger, retry)
		gotPortfolio, err := s.getPortfolio()
		if err != nil {
			t.Fatalf("getPortfolio вернул ошибку: %v", err)
		}

		core_testutils.AssertEqual(t, gotPortfolio, mock.portfolio)
	})

	t.Run("клиент вернул ошибку", func(t *testing.T) {
		mock := &mockClient{
			portfolio: &Portfolio{
				TotalAmount: -1,
			},
			err: core_errors.New(
				core_errors.ErrExternalAPI,
				"api down",
			),
		}
		s := newService(mock, logger, retry)
		_, err := s.getPortfolio()
		if err == nil {
			t.Fatal("ожидали ошибку, а получили nil")
		}

		if !errors.Is(err, mock.err) {
			t.Fatalf("ожидали ErrExternalAPI, получили; %v", err)
		}
	})
}
