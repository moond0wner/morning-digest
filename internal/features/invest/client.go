package invest

import (
	"context"
	"fmt"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"opensource.tbank.ru/invest/invest-go/investgo"
	investapi "opensource.tbank.ru/invest/invest-go/proto"
)

type client struct {
	api               *investgo.Client
	usersService      *investgo.UsersServiceClient
	operationsService *investgo.OperationsServiceClient
	accountID         string
}

func newClient(ctx context.Context, logger *core_logger.Logger) (*client, error) {
	config := newConfigMust()

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
		return nil, fmt.Errorf("create invest client: %w", err)
	}

	return &client{
		api:               api,
		usersService:      api.NewUsersServiceClient(),
		operationsService: api.NewOperationsServiceClient(),
		accountID:         config.AccountId,
	}, nil
}

func (c *client) GetPortfolio(ctx context.Context) (*investgo.PortfolioResponse, error) {
	response, err := c.operationsService.GetPortfolio(
		c.accountID,
		investapi.PortfolioRequest_RUB,
	)
	if err != nil {
		return nil, fmt.Errorf("get portfolio: %w", err)
	}
	return response, nil
}
