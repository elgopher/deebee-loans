// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package database

import (
	"sync"

	"github.com/jacekolszak/deebee-loans/service"
)

// SynchronizedLoans serialize service.Loans method calls, so the `service` package is free of low-level code such as synchronization.
// This approach is simple, but it might be inefficient for large states (more than 100MB of memory), because some methods
// might block access to state for a significant amount of time (for example SynchronizedLoans.Snapshot may take seconds).
type SynchronizedLoans struct {
	mutex sync.Mutex
	loans *service.Loans
}

func (s *SynchronizedLoans) TakeLoan(userID string, amount, term int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.loans.TakeLoan(userID, amount, term)
}

func (s *SynchronizedLoans) PayLoan(userID string, amount int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.loans.PayLoan(userID, amount)
}

func (s *SynchronizedLoans) GetActiveLoan(userID string) (service.ActiveLoan, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.loans.GetActiveLoan(userID)
}

func (s *SynchronizedLoans) Snapshot() service.Snapshot {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.loans.Snapshot()
}
