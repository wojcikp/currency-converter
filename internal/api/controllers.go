package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (s *GinServer) GetRates(c *gin.Context) {
	param := c.Query("currencies")
	if param == "" {
		logrus.Error(`url parameter "currencies" not provided`)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	validatedCurrencies, err := validateCurrencies(strings.Split(param, ","))
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return

	}

	rates, err := s.converter.GetCurrenciesRates(c.Request.Context(), validatedCurrencies)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, rates)
}

func (s *GinServer) ExchangeCryptoCurrencies(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amount := c.Query("amount")

	if from == "" || to == "" || amount == "" {
		logrus.Errorf("missing one of parameters: from, to or amount. parameters: from: %s, to: %s, amount: %s", from, to, amount)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	decimalAmount, err := decimal.NewFromString(amount)
	if err != nil {
		logrus.Error("could not parse parameter amount to decimal. amount: ", amount)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	exchangedCrypto, err := s.converter.ConvertCryptoCurrencies(
		c.Request.Context(),
		strings.ToUpper(from),
		strings.ToUpper(to),
		decimalAmount,
	)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, exchangedCrypto)
}

func validateCurrencies(currencies []string) ([]string, error) {
	currenciesSet := map[string]struct{}{}
	var validatedCurrencies []string
	for _, currency := range currencies {
		currency = strings.ToUpper(strings.TrimSpace(currency))
		if currency == "" {
			continue
		}
		if _, ok := currenciesSet[currency]; !ok {
			currenciesSet[currency] = struct{}{}
		}
	}
	if len(currenciesSet) < 2 {
		return nil, fmt.Errorf("not enough currencies to exchange provided, currencies: %s", currencies)
	}
	for currency := range currenciesSet {
		validatedCurrencies = append(validatedCurrencies, currency)
	}
	return validatedCurrencies, nil
}
