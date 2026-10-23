package invest

import (
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
)

const rubTicker = "RUB000UTSTOM"

type investClient interface {
	getPortfolio() (*Portfolio, error)
	close() error
}

type service struct {
	client investClient
	log    *core_logger.Logger
	retry  *core_retry.Retrier
}

func newService(client investClient, logger *core_logger.Logger, retrier *core_retry.Retrier) *service {
	return &service{
		client: client,
		log:    logger,
		retry:  retrier,
	}
}

func (s *service) close() error {
	return s.client.close()
}

func (s *service) getPortfolio() (*Portfolio, error) {
	portfolio, err := s.client.getPortfolio()
	if err != nil {
		return nil, err
	}

	return portfolio, nil
}
