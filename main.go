// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/jacekolszak/deebee-loans/database"
	"github.com/jacekolszak/deebee-loans/web"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	loans, done, err := database.Start(ctx)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	defer func() {
		<-done // wait until database saves snapshot of service.Loans
	}()

	if err := web.ListenAndServe(ctx, loans); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Error starting server")
	}
}
