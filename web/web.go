// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jacekolszak/deebee-loans/service"
	"github.com/jacekolszak/yala/logger"
)

func ListenAndServe(ctx context.Context, loans Loans) error {
	mux := http.NewServeMux()
	mux.Handle("/take", takeLoan{loans: loans})
	mux.Handle("/pay", payLoan{loans: loans})

	server := &http.Server{Addr: ":8080", Handler: mux}
	shutdownServerOnceDone(ctx, server)

	logger.With(ctx, "address", server.Addr).Info("Starting web server")
	return server.ListenAndServe()
}

type Loans interface {
	TakeLoan(userID string, amount, term int) error
	PayLoan(userID string, amount int) error
	GetActiveLoan(userID string) (service.ActiveLoan, error)
}

func shutdownServerOnceDone(ctx context.Context, server *http.Server) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info(ctx, "Shutting down web server")
				if err := server.Shutdown(context.Background()); err != nil {
					logger.WithError(ctx, err).Error("Problem shutting down the server")
				}
				return
			}
		}
	}()
}

func writeClientError(writer http.ResponseWriter, o ...interface{}) {
	writer.WriteHeader(400)
	_, _ = fmt.Fprintln(writer, o...)
}
