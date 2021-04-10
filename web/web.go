// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jacekolszak/deebee-loans/service"
	"github.com/sirupsen/logrus"
)

type Loans interface {
	TakeLoan(userID string, amount, term int) error
	PayLoan(userID string, amount int) error
	GetActiveLoan(userID string) (service.ActiveLoan, error)
}

func ListenAndServe(ctx context.Context, loans Loans) error {
	mux := http.NewServeMux()
	mux.Handle("/take", takeLoan{loans: loans})
	mux.Handle("/pay", payLoan{loans: loans})

	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		for {
			select {
			case <-ctx.Done():
				logrus.Info("Shutting down web server")
				if err := server.Shutdown(context.Background()); err != nil {
					logrus.WithError(err).Error("Problem shutting down the server")
				}
				return
			}
		}
	}()

	logrus.Infof("Starting web server on %s", server.Addr)
	return server.ListenAndServe()
}

type takeLoan struct {
	loans Loans
}

func (h takeLoan) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		writeClientError(writer, err)
		return
	}

	var (
		userID    = request.Form.Get("user")
		term, _   = strconv.Atoi(request.Form.Get("term"))
		amount, _ = strconv.Atoi(request.Form.Get("amount"))
	)
	err := h.loans.TakeLoan(userID, term, amount)
	if err != nil {
		writeClientError(writer, err)
		return
	}
	writer.WriteHeader(201)
}

type payLoan struct {
	loans Loans
}

func (h payLoan) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		writeClientError(writer, err)
		return
	}

	var (
		userID    = request.Form.Get("user")
		amount, _ = strconv.Atoi(request.Form.Get("amount"))
	)
	err := h.loans.PayLoan(userID, amount)
	if err != nil {
		writeClientError(writer, err)
		return
	}
	activeLoan, err := h.loans.GetActiveLoan(userID)
	if err != nil {
		writer.WriteHeader(500)
		logrus.WithError(err).Error("error getting active loan")
		return
	}
	_, _ = fmt.Fprintln(writer, "Amount remaining", activeLoan.AmountRemaining)
}

func writeClientError(writer http.ResponseWriter, o ... interface{}) {
	writer.WriteHeader(400)
	_, _ = fmt.Fprintln(writer, o...)
}