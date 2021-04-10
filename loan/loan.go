// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package loan

import (
	"errors"
	"fmt"
)

type Loan struct {
	amount int
	term   int
	paid   int
}

func New(amount, term int) *Loan {
	return &Loan{
		amount: amount,
		term:   term,
	}
}

func FromSnapshot(snapshot Snapshot) *Loan {
	return &Loan{
		amount: snapshot.Amount,
		term:   snapshot.Term,
		paid:   snapshot.Paid,
	}
}

type Snapshot struct {
	Amount int
	Term   int
	Paid   int
}

func (l *Loan) Snapshot() Snapshot {
	return Snapshot{
		Amount: l.amount,
		Term:   l.term,
		Paid:   l.paid,
	}
}

func (l *Loan) Pay(amount int) error {
	if l.paid+amount > l.amount {
		return fmt.Errorf("too much, maximum amount to pay is %d", l.AmountRemaining())
	}
	if amount < 0 {
		return errors.New("negative amount")
	}
	l.paid += amount
	return nil
}

func (l *Loan) AmountRemaining() int {
	return l.amount - l.paid
}
