// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package web

import (
	"net/http"
	"strconv"
)

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
