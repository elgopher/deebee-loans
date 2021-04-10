// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package service

import (
	"fmt"

	"github.com/jacekolszak/deebee-loans/loan"
)

type Loans struct {
	activeLoansByUser map[string]*loan.Loan
}

func FromSnapshot(snapshot Snapshot) *Loans {
	loans := &Loans{
		activeLoansByUser: map[string]*loan.Loan{},
	}
	for userID, activeLoan := range snapshot.Loans {
		loans.activeLoansByUser[userID] = loan.FromSnapshot(activeLoan)
	}
	return loans
}

type Snapshot struct {
	Loans map[string]loan.Snapshot
}

func (l *Loans) Snapshot() Snapshot {
	snapshot := Snapshot{
		Loans: map[string]loan.Snapshot{},
	}
	for userID, activeLoan := range l.activeLoansByUser {
		snapshot.Loans[userID] = activeLoan.Snapshot()
	}
	return snapshot
}

func (l *Loans) TakeLoan(userID string, amount, term int) error {
	_, alreadyHasLoan := l.activeLoansByUser[userID]
	if alreadyHasLoan {
		return fmt.Errorf("user %s already has active loan", userID)
	}

	newLoan := loan.New(amount, term)
	l.activeLoansByUser[userID] = newLoan
	return nil
}

func (l *Loans) PayLoan(userID string, amount int) error {
	userLoan, found := l.activeLoansByUser[userID]
	if !found {
		return fmt.Errorf("no active loan found for user %s", userID)
	}
	return userLoan.Pay(amount)
}

func (l *Loans) GetActiveLoan(userID string) (ActiveLoan, error) {
	activeLoan, found := l.activeLoansByUser[userID]
	if !found {
		return ActiveLoan{}, fmt.Errorf("no active loan found for user %s", userID)
	}
	return ActiveLoan{
		AmountRemaining: activeLoan.AmountRemaining(),
	}, nil
}

type ActiveLoan struct {
	AmountRemaining int
}
