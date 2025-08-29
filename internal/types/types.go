package types

import (
	"context"

	"github.com/shopspring/decimal"
)

type RatesProvider interface {
	GetExchangeRates(context.Context) (map[string]decimal.Decimal, error)
	GetCryptoExchangeRates(ctx context.Context) map[string]CryptoCurrencyInfo
}

type Converter interface {
	GetCurrenciesRates(ctx context.Context, currencies []string) ([]ConvertedRate, error)
	ConvertCryptoCurrencies(ctx context.Context, from, to string, amount decimal.Decimal) (ExchangedCryptoCurrency, error)
}

type ConvertedRate struct {
	From string          `json:"from"`
	To   string          `json:"to"`
	Rate decimal.Decimal `json:"rate"`
}

type ExchangedCryptoCurrency struct {
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}

type CryptoCurrencyInfo struct {
	DecimalPlaces int
	RateToUSD     decimal.Decimal
}
