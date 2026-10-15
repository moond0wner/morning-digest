package invest

import (
	"os"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
)

const (
	prodAddress    = "invest-public-api.tbank.ru:443"
	sandboxAddress = "sandbox-invest-public-api.tbank.ru:443"
)

type config struct {
	// EndPoint - Для работы с реальным контуром и контуром песочницы нужны разные эндпоинты.
	// По умолчанию = sandbox-invest-public-api.tbank.ru:443
	// https://tinkoff.github.io/investAPI/url_difference/
	EndPoint string `yaml:"EndPoint"`
	// Token - Ваш токен для InvestAPI
	Token string `yaml:"APIToken"`
	// AppName - Название вашего приложения, по умолчанию = tinkoff-api-go-sdk
	AppName string `yaml:"AppName"`
	// AccountId - Если уже есть аккаунт для апи можно указать напрямую,
	// по умолчанию откроется новый счет в песочнице
	AccountId string `yaml:"AccountId"`
	// DisableResourceExhaustedRetry - Если true, то сдк не пытается ретраить, после получения ошибки об исчерпывании
	// лимита запросов, если false, то сдк ждет нужное время и пытается выполнить запрос снова. По умолчанию = false
	DisableResourceExhaustedRetry bool `yaml:"DisableResourceExhaustedRetry"`
	// DisableAllRetry - Отключение всех ретраев
	DisableAllRetry bool `yaml:"DisableAllRetry"`
	// MaxRetries - Максимальное количество попыток переподключения, по умолчанию = 3
	// (если указать значение 0 это не отключит ретраи, для отключения нужно прописать DisableAllRetry = true)
	MaxRetries uint `yaml:"MaxRetries"`
	// InsecureSkipVerify - Пропустить проверку сертификата сервера
	InsecureSkipVerify bool `yaml:"InsecureSkipVerify"` // никаких сертификатов
}

func newConfig() (config, error) {
	token := os.Getenv("INVEST_TOKEN")
	if token == "" {
		return config{}, core_errors.New(core_errors.ErrInternal, "empty INVEST_TOKEN")
	}
	accountId := os.Getenv("INVEST_ACCOUNT_ID")
	if accountId == "" {
		return config{}, core_errors.New(core_errors.ErrInternal, "empty INVEST_ACCOUNT_ID")
	}

	return config{
		EndPoint:                      prodAddress,
		Token:                         token,
		AppName:                       "morning-digest",
		AccountId:                     accountId,
		DisableResourceExhaustedRetry: false,
		DisableAllRetry:               false,
		MaxRetries:                    3,
		InsecureSkipVerify:            true,
	}, nil
}
