package currencyconverter

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	exchangeratesprovider "github.com/wojcikp/currency-converter/internal/exchange_rates_provider"
	"github.com/wojcikp/currency-converter/internal/types"
)

func TestGetCurrenciesRates(t *testing.T) {
	provider := exchangeratesprovider.NewExchangeRatesProviderMock()
	converter := NewConverter(provider)
	ctx := context.Background()
	cases := []struct {
		name      string
		in        []string
		wantCount int
	}{
		{"three currencies", []string{"USD", "GBP", "EUR"}, 6},
		{"two currencies", []string{"GBP", "EUR"}, 2},
		{"one currency", []string{"GBP"}, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := converter.GetCurrenciesRates(ctx, tc.in)
			if err != nil {
				t.Fatalf("TestGetCurrenciesRates error: %v", err)
			}
			if len(out) != tc.wantCount {
				t.Fatalf("len=%d, want %d", len(out), tc.wantCount)
			}
		})
	}
}

func TestGetCurrencyPairsToExchange(t *testing.T) {
	cases := []struct {
		name       string
		currencies []string
		expected   []types.ConvertedRate
	}{
		{
			name:       "three currencies",
			currencies: []string{"USD", "EUR", "GBP"},
			expected: []types.ConvertedRate{
				{From: "USD", To: "EUR", Rate: decimal.Decimal{}},
				{From: "EUR", To: "USD", Rate: decimal.Decimal{}},
				{From: "USD", To: "GBP", Rate: decimal.Decimal{}},
				{From: "GBP", To: "USD", Rate: decimal.Decimal{}},
				{From: "EUR", To: "GBP", Rate: decimal.Decimal{}},
				{From: "GBP", To: "EUR", Rate: decimal.Decimal{}},
			},
		},
		{
			name:       "two currencies",
			currencies: []string{"USD", "EUR"},
			expected: []types.ConvertedRate{
				{From: "USD", To: "EUR", Rate: decimal.Decimal{}},
				{From: "EUR", To: "USD", Rate: decimal.Decimal{}},
			},
		},
		{
			name:       "one currency",
			currencies: []string{"USD"},
			expected:   []types.ConvertedRate{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := getCurrencyPairsToExchange(tc.currencies)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}
