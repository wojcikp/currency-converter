package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
	currencyconverter "github.com/wojcikp/currency-converter/internal/currency_converter"
	exchangeratesprovider "github.com/wojcikp/currency-converter/internal/exchange_rates_provider"
	"github.com/wojcikp/currency-converter/internal/types"
)

func setupRouter() *gin.Engine {
	provider := exchangeratesprovider.NewExchangeRatesProviderMock()
	converter := currencyconverter.NewConverter(provider)
	server := NewGinServer("8080", converter)
	router := gin.Default()
	router.GET("/rates", server.GetRates)
	router.GET("/exchange", server.ExchangeCryptoCurrencies)
	return router
}

func TestRatesEndpoint(t *testing.T) {
	router := setupRouter()

	cases := []struct {
		name              string
		url               string
		wantStatus        int
		wantRates         []types.ConvertedRate
		wantEmptyResponse bool
	}{
		{
			name:       "three currencies",
			url:        "/rates?currencies=USD,GBP,EUR",
			wantStatus: 200,
			wantRates: []types.ConvertedRate{
				{From: "USD", To: "GBP", Rate: decimal.RequireFromString("0.743283")},
				{From: "GBP", To: "USD", Rate: decimal.RequireFromString("1.3453825797172813")},
				{From: "USD", To: "EUR", Rate: decimal.RequireFromString("0.861355")},
				{From: "EUR", To: "USD", Rate: decimal.RequireFromString("1.1609615083211916")},
				{From: "GBP", To: "EUR", Rate: decimal.RequireFromString("1.1588520119523788")},
				{From: "EUR", To: "GBP", Rate: decimal.RequireFromString("0.8629229527895003")},
			},
		},
		{
			name:       "two currencies",
			url:        "/rates?currencies=GBP,EUR",
			wantStatus: 200,
			wantRates: []types.ConvertedRate{
				{From: "GBP", To: "EUR", Rate: decimal.RequireFromString("1.1588520119523788")},
				{From: "EUR", To: "GBP", Rate: decimal.RequireFromString("0.8629229527895003")},
			},
		},
		{
			name:              "one currency",
			url:               "/rates?currencies=GBP",
			wantStatus:        400,
			wantRates:         []types.ConvertedRate{},
			wantEmptyResponse: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("TestRatesEndpoint error: %v", err)
			}

			router.ServeHTTP(w, req)
			if w.Code != tc.wantStatus {
				t.Fatalf("response status=%d, want %d", w.Code, tc.wantStatus)
			}

			if tc.wantEmptyResponse {
				var m gin.H
				if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
					t.Fatalf("expected valid JSON error: %v", err)
				}
				return
			}

			var got []types.ConvertedRate
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("cannot unmarshal: %v", err)
			}

			if len(got) != len(tc.wantRates) {
				t.Fatalf("unexpected number of rates: got=%d, want=%d", len(got), len(tc.wantRates))
			}

			sortRates(got)
			sortRates(tc.wantRates)

			if diff := cmp.Diff(tc.wantRates, got); diff != "" {
				t.Errorf("%s test mismatch (-want +got):\n%s", tc.name, diff)
			}
		})
	}
}
func TestExchangeEndpoint(t *testing.T) {
	router := setupRouter()

	cases := []struct {
		name              string
		url               string
		wantStatus        int
		wantResponse      types.ExchangedCryptoCurrency
		wantEmptyResponse bool
	}{
		{
			name:       "WBTC to USDT",
			url:        "/exchange?from=WBTC&to=USDT&amount=1.0",
			wantStatus: 200,
			wantResponse: types.ExchangedCryptoCurrency{
				From:   "WBTC",
				To:     "USDT",
				Amount: decimal.RequireFromString("57094.314314"),
			},
		},
		{
			name:       "USDT to BEER",
			url:        "/exchange?from=USDT&to=BEER&amount=1.0",
			wantStatus: 200,
			wantResponse: types.ExchangedCryptoCurrency{
				From:   "USDT",
				To:     "BEER",
				Amount: decimal.RequireFromString("40593.2547744819179195"),
			},
		},
		{
			name:              "MATIC to GATE",
			url:               "/exchange?from=MATIC&to=GATE&amount=0.999",
			wantStatus:        400,
			wantResponse:      types.ExchangedCryptoCurrency{},
			wantEmptyResponse: true,
		},
		{
			name:              "USDT to GATE, no amount",
			url:               "/exchange?from=USDT&to=GATE",
			wantStatus:        400,
			wantResponse:      types.ExchangedCryptoCurrency{},
			wantEmptyResponse: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("TestExchangeEndpoint error: %v", err)
			}

			router.ServeHTTP(w, req)
			if w.Code != tc.wantStatus {
				t.Fatalf("response status=%d, want %d", w.Code, tc.wantStatus)
			}

			if tc.wantEmptyResponse {
				var m gin.H
				if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
					t.Fatalf("expected valid JSON error: %v", err)
				}
				return
			}

			var got types.ExchangedCryptoCurrency
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("cannot unmarshal: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, got); diff != "" {
				t.Errorf("%s test mismatch (-want +got):\n%s", tc.name, diff)
			}
		})
	}
}

func sortRates(rates []types.ConvertedRate) {
	sort.Slice(rates, func(i, j int) bool {
		if rates[i].From == rates[j].From {
			return rates[i].To < rates[j].To
		}
		return rates[i].From < rates[j].From
	})
}
