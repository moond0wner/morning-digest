package invest

import (
	"context"
	"fmt"

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

func (c *client) GetPortfolio() (*investgo.PortfolioResponse, error) {
	response, err := c.operationsService.GetPortfolio(
		c.accountID,
		investapi.PortfolioRequest_RUB,
	)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrExternalAPI,
			err,
			"get portfolio",
		)
	}
	return response, nil
}

func (c *client) Close() error {
	if c.api != nil {
		return c.api.Conn.Close()
	}
	return nil
}
