package exchangeratesprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wojcikp/currency-converter/internal/types"
)

type ExchangeRatesProvider struct {
	openexchangeAppId string
	httpClient        *http.Client
}
type ExchangeRatesProviderMock struct{}

type ExchangeRates struct {
	Rates map[string]decimal.Decimal `json:"rates"`
}

func NewExchangeRatesProvider(openexchangeAppId string) (*ExchangeRatesProvider, error) {
	return &ExchangeRatesProvider{
		openexchangeAppId: openexchangeAppId,
		httpClient:        &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (p *ExchangeRatesProvider) GetExchangeRates(ctx context.Context) (map[string]decimal.Decimal, error) {
	url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", p.openexchangeAppId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request err: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during openexchangerates.org api GET, err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var data ExchangeRates
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("json decoding error: %w", err)
	}

	return data.Rates, nil
}

func (p *ExchangeRatesProvider) GetCryptoExchangeRates(ctx context.Context) map[string]types.CryptoCurrencyInfo {
	return map[string]types.CryptoCurrencyInfo{
		"BEER":  {DecimalPlaces: 18, RateToUSD: decimal.RequireFromString("0.00002461")},
		"FLOKI": {DecimalPlaces: 18, RateToUSD: decimal.RequireFromString("0.0001428")},
		"GATE":  {DecimalPlaces: 18, RateToUSD: decimal.RequireFromString("6.87")},
		"USDT":  {DecimalPlaces: 6, RateToUSD: decimal.RequireFromString("0.999")},
		"WBTC":  {DecimalPlaces: 8, RateToUSD: decimal.RequireFromString("57037.22")},
	}
}

func NewExchangeRatesProviderMock() *ExchangeRatesProviderMock {
	return &ExchangeRatesProviderMock{}
}

func (p *ExchangeRatesProviderMock) GetExchangeRates(ctx context.Context) (map[string]decimal.Decimal, error) {
	rates := make(map[string]decimal.Decimal)
	rates["EUR"] = decimal.NewFromFloat(0.861355)
	rates["GBP"] = decimal.NewFromFloat(0.743283)
	rates["USD"] = decimal.NewFromInt(1)
	return rates, nil
}

func (p *ExchangeRatesProviderMock) GetCryptoExchangeRates(ctx context.Context) map[string]types.CryptoCurrencyInfo {
	return map[string]types.CryptoCurrencyInfo{
		"BEER":  {DecimalPlaces: 18, RateToUSD: decimal.RequireFromString("0.00002461")},
		"FLOKI": {DecimalPlaces: 18, RateToUSD: decimal.RequireFromString("0.0001428")},
		"GATE":  {DecimalPlaces: 18, RateToUSD: decimal.RequireFromString("6.87")},
		"USDT":  {DecimalPlaces: 6, RateToUSD: decimal.RequireFromString("0.999")},
		"WBTC":  {DecimalPlaces: 8, RateToUSD: decimal.RequireFromString("57037.22")},
	}
}
