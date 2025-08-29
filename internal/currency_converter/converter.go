package currencyconverter

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/wojcikp/currency-converter/internal/types"
)

type Converter struct {
	exchangeRatesProvider types.RatesProvider
}

func NewConverter(ratesProvider types.RatesProvider) *Converter {
	return &Converter{ratesProvider}
}

func (c *Converter) GetCurrenciesRates(ctx context.Context, currencies []string) ([]types.ConvertedRate, error) {
	rates, err := c.exchangeRatesProvider.GetExchangeRates(ctx)
	if err != nil {
		return []types.ConvertedRate{}, fmt.Errorf("error during fetching exchange rates, err: %w", err)
	}

	if err = validateCurrencies(currencies, rates); err != nil {
		return []types.ConvertedRate{}, err
	}

	exchangePairs := getCurrencyPairsToExchange(currencies)

	for i := range exchangePairs {
		exchangePairs[i].Rate = rates[exchangePairs[i].To].Div(rates[exchangePairs[i].From])
	}

	return exchangePairs, nil
}

func (c *Converter) ConvertCryptoCurrencies(
	ctx context.Context,
	from, to string,
	amount decimal.Decimal,
) (types.ExchangedCryptoCurrency, error) {
	rates := c.exchangeRatesProvider.GetCryptoExchangeRates(ctx)

	currencyFrom, ok := rates[from]
	if !ok {
		return types.ExchangedCryptoCurrency{}, fmt.Errorf("currency: %s not found in crypto currency rates", from)
	}
	currencyTo, ok := rates[to]
	if !ok {
		return types.ExchangedCryptoCurrency{}, fmt.Errorf("currency: %s not found in crypto currency rates", to)
	}

	usd := amount.Mul(currencyFrom.RateToUSD)
	result := usd.Div(currencyTo.RateToUSD)
	result = result.Round(int32(currencyTo.DecimalPlaces))

	return types.ExchangedCryptoCurrency{From: from, To: to, Amount: result}, nil
}

func validateCurrencies(currencies []string, rates map[string]decimal.Decimal) error {
	for _, currency := range currencies {
		_, ok := rates[currency]
		if !ok {
			return fmt.Errorf("currency: %s not found in openexchangerates.org rates", currency)
		}
	}
	return nil
}

func getCurrencyPairsToExchange(currencies []string) []types.ConvertedRate {
	var currencyPairs []types.ConvertedRate

	for i := 0; i < len(currencies); i++ {
		for j := i + 1; j < len(currencies); j++ {
			currencyPairs = append(currencyPairs, types.ConvertedRate{From: currencies[i], To: currencies[j]})
			currencyPairs = append(currencyPairs, types.ConvertedRate{From: currencies[j], To: currencies[i]})
		}
	}

	return currencyPairs
}
