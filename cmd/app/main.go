package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/wojcikp/currency-converter/internal/app"
)

func main() {
	app, err := app.BuildApp()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Application built successfully. Running app...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go app.Run()

	<-ctx.Done()
	logrus.Info("Shutdown signal received, shutting down application...")

	if err := app.Shutdown(); err != nil {
		logrus.Fatal("Application forced to shutdown: ", err)
	}
}
