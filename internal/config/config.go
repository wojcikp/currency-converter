package config

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerPort        string
	OpenExchangeAppID string
}

func Load() (*Config, error) {
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		logrus.Warn("could not read SERVER_PORT env variable. Setting server port to default :8080")
		serverPort = "8080"
	}
	openexchangeAppId := os.Getenv("OPENEXCHANGE_APP_ID")
	if openexchangeAppId == "" {
		return nil, errors.New("could not read OPENEXCHANGE_APP_ID env variable. provide OPENEXCHANGE_APP_ID env variable to run application")
	}
	return &Config{serverPort, openexchangeAppId}, nil
}
