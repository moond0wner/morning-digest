package invest

import (
	"context"
	"fmt"
	"sort"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"opensource.tbank.ru/invest/invest-go/investgo"
	investapi "opensource.tbank.ru/invest/invest-go/proto"
)

type client struct {
	api               *investgo.Client
	operationsService *investgo.OperationsServiceClient
	accountID         string
}

func newClient(ctx context.Context, logger *core_logger.Logger) (*client, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("load invest config: %w", err)
	}

	investgoConfig := investgo.Config{
		EndPoint:                      config.EndPoint,
		Token:                         config.Token,
		AppName:                       config.AppName,
		AccountId:                     config.AccountId,
		DisableResourceExhaustedRetry: config.DisableResourceExhaustedRetry,
		DisableAllRetry:               config.DisableAllRetry,
		MaxRetries:                    config.MaxRetries,
		InsecureSkipVerify:            config.InsecureSkipVerify,
	}

	api, err := investgo.NewClient(ctx, investgoConfig, logger.Sugar())
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrExternalAPI,
			err,
			"create invest client",
		)
	}

	return &client{
		api:               api,
		operationsService: api.NewOperationsServiceClient(),
		accountID:         config.AccountId,
	}, nil
}

func (c *client) getPortfolio(ctx context.Context) (*Portfolio, error) {
	type result struct {
		p   *Portfolio
		err error
	}
	resCh := make(chan result, 1)
	go func() {
		response, err := c.operationsService.GetPortfolio(
			c.accountID,
			investapi.PortfolioRequest_RUB,
		)
		if err != nil {
			resCh <- result{err: core_errors.Wrap(
				core_errors.ErrExternalAPI,
				err,
				"get portfolio",
			)}
			return
		}
		resCh <- result{p: convertToPortfolio(response)}
	}()
	select {
	case <-ctx.Done():
		return nil, core_errors.Wrap(
			core_errors.ErrExternalAPI,
			ctx.Err(),
			"get portfolio timeout",
		)
	case res := <-resCh:
		return res.p, res.err
	}

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

func (c *client) close() error {
	if c.api != nil {
		return c.api.Conn.Close()
	}
	return nil
}
