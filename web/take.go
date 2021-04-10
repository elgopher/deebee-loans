// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package web

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type takeLoan struct {
	loans Loans
}

func (h takeLoan) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var (
		userID    = request.FormValue("user")
		term, _   = strconv.Atoi(request.FormValue("term"))
		amount, _ = strconv.Atoi(request.FormValue("amount"))
	)

	err := h.loans.TakeLoan(userID, term, amount)
	if err != nil {
		writeClientError(writer, err)
		return
	}

	logrus.Info("New loan taken")
	writer.WriteHeader(201)
}
