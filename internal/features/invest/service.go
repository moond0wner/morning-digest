package invest

import (
	"context"
	"fmt"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"opensource.tbank.ru/invest/invest-go/investgo"
)

type Service struct {
	client *client
}

func NewService(ctx context.Context, logger *core_logger.Logger) *Service {
	clnt, err := newClient(ctx, logger)
	if err != nil {
		panic(err)
	}
	return &Service{
		client: clnt,
	}
}

func (s *Service) GetPortfolio(ctx context.Context) (*Portfolio, error) {
	resp, err := s.client.GetPortfolio(ctx)
	if err != nil {
		return nil, fmt.Errorf("get portfolio: %w", err)
	}

	return convertToPortfolio(resp), nil
}

func convertToPortfolio(resp *investgo.PortfolioResponse) *Portfolio {
	p := &Portfolio{
		TotalAmount:   resp.GetTotalAmountPortfolio().ToFloat(),
		ExpectedYield: resp.GetExpectedYield().ToFloat(),
		DailyYield:    resp.GetDailyYield().ToFloat(),
	}

	for _, pos := range resp.GetPositions() {
		if pos.Ticker == "RUB000UTSTOM" { // денежный остаток - не нужен
			continue
		}

		p.Positions = append(p.Positions, Position{
			Ticker:        pos.GetTicker(),
			Quantity:      pos.GetQuantity().ToFloat(),
			CurrentPrice:  pos.GetCurrentPrice().ToFloat(),
			ExpectedYield: pos.GetExpectedYield().ToFloat(),
		})
	}

	return p
}
