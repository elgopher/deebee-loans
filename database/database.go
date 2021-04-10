// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jacekolszak/deebee-loans/service"
	"github.com/jacekolszak/deebee/codec"
	"github.com/jacekolszak/deebee/json"
	"github.com/jacekolszak/deebee/store"
	"github.com/sirupsen/logrus"
)

func StartLoans(ctx context.Context) (loans *SynchronizedLoans, done <-chan struct{}, err error) {
	s, err := store.Open("/tmp/loans")
	if err != nil {
		return nil, nil, fmt.Errorf("error opening DeeBee store: %w", err)
	}

	loans, err = loadState(s)
	if err != nil {
		return nil, nil, err
	}

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		for {
			select {
			case <-time.After(time.Minute):
				saveState(loans, s)
			case <-ctx.Done():
				saveState(loans, s)
				return
			}
		}
	}()
	return loans, doneCh, nil
}

func loadState(s *store.Store) (*SynchronizedLoans, error) {
	snapshot := service.Snapshot{}
	_, err := codec.ReadLatest(s, json.Decoder(&snapshot))
	if err != nil && !store.IsVersionNotFound(err) {
		return nil, err
	}
	loans := service.FromSnapshot(snapshot)

	return &SynchronizedLoans{
		loans: loans,
	}, nil
}

func saveState(loans *SynchronizedLoans, s *store.Store) {
	logrus.Info("Saving loans.Service state")
	snapshot := loans.Snapshot()
	if err := json.Write(s, snapshot); err != nil {
		logrus.WithError(err).Error("error saving state")
	}
}
