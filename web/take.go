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
		amount, _ = strconv.Atoi(request.FormValue("amount"))
		term, _   = strconv.Atoi(request.FormValue("term"))
	)

	err := h.loans.TakeLoan(userID, amount, term)
	if err != nil {
		writeClientError(writer, err)
		return
	}

	logrus.Info("New loan taken")
	writer.WriteHeader(201)
}
