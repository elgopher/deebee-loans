// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type payLoan struct {
	loans Loans
}

func (h payLoan) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var (
		userID    = request.FormValue("user")
		amount, _ = strconv.Atoi(request.FormValue("amount"))
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

	logrus.Info("Loan paid off")
	_, _ = fmt.Fprintln(writer, "Amount remaining:", activeLoan.AmountRemaining)
}
