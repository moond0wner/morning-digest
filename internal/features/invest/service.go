package invest

import (
	"context"
	"fmt"
	"sort"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

const rubTicker = "RUB000UTSTOM"

type Service struct {
	client *client
}

func NewService(ctx context.Context, logger *core_logger.Logger) (*Service, error) {
	clnt, err := newClient(ctx, logger)
	if err != nil {
		return nil, fmt.Errorf("create invest service: %w", err)
	}
	return &Service{
		client: clnt,
	}, nil
}

func (s *Service) Close() error {
	return s.client.Close()
}

func (s *Service) GetPortfolio() (*Portfolio, error) {
	response, err := s.client.GetPortfolio()
	if err != nil {
		return nil, fmt.Errorf("get portfolio: %w", err)
	}

	return convertToPortfolio(response), nil
}

func convertToPortfolio(resp *investgo.PortfolioResponse) *Portfolio {
	p := &Portfolio{
		TotalAmount:   resp.GetTotalAmountPortfolio().ToFloat(),
		ExpectedYield: resp.GetExpectedYield().ToFloat(),
		DailyYield:    resp.GetDailyYield().ToFloat(),
	}

	for _, pos := range resp.GetPositions() {
		if pos.Ticker == rubTicker { // денежный остаток - не нужен
			continue
		}

		p.Positions = append(p.Positions, Position{
			Ticker:        pos.GetTicker(),
			Quantity:      pos.GetQuantity().ToFloat(),
			CurrentPrice:  pos.GetCurrentPrice().ToFloat(),
			ExpectedYield: pos.GetExpectedYield().ToFloat(),
		})
	}
	sort.Slice(p.Positions, func(i, j int) bool {
		return p.Positions[i].CurrentPrice*p.Positions[i].Quantity >
			p.Positions[j].CurrentPrice*p.Positions[j].Quantity
	})
	return p
}
