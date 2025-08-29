package app

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wojcikp/currency-converter/internal/api"
	"github.com/wojcikp/currency-converter/internal/config"
	currencyconverter "github.com/wojcikp/currency-converter/internal/currency_converter"
	exchangeratesprovider "github.com/wojcikp/currency-converter/internal/exchange_rates_provider"
)

type App struct {
	server *api.GinServer
}

func BuildApp() (*App, error) {
	config, err := config.Load()
	if err != nil {
		return nil, err
	}
	logrus.Info("Application config loaded successfully")

	ratesProvider, err := exchangeratesprovider.NewExchangeRatesProvider(config.OpenExchangeAppID)
	if err != nil {
		return nil, err
	}
	logrus.Info("Rates provider initialized")

	converter := currencyconverter.NewConverter(ratesProvider)
	logrus.Info("Currency converter initialized")

	server := api.NewGinServer(config.ServerPort, converter)
	logrus.Info("Gin server initialized")

	return &App{server}, nil
}

func (a *App) Run() {
	a.server.RegisterRoutes()
	if err := a.server.Run(); err != nil && err != http.ErrServerClosed {
		logrus.Fatal("Could not run the application server due to an error:", err)
	}
}

func (a *App) Shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	logrus.Info("Server is off")
	return nil
}
