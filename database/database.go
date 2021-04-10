// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package database

import (
	"context"
	"time"

	"github.com/jacekolszak/deebee-loans/service"
	"github.com/jacekolszak/deebee/store"
	"github.com/sirupsen/logrus"
)

func StartLoans(ctx context.Context, s Store) (loans *SynchronizedLoans, done <-chan struct{}, err error) {
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

type Store interface {
	ReadLatest(out *service.Snapshot) (store.Version, error)
	Write(in *service.Snapshot, options ...store.WriterOption) error
}

func loadState(s Store) (*SynchronizedLoans, error) {
	snapshot := service.Snapshot{}
	version, err := s.ReadLatest(&snapshot)
	if err != nil && !store.IsVersionNotFound(err) {
		return nil, err
	}

	logrus.Infof("Snapshot loaded with version %+v", version)
	loans := service.FromSnapshot(snapshot)

	return &SynchronizedLoans{
		loans: loans,
	}, nil
}

func saveState(loans *SynchronizedLoans, s Store) {
	logrus.Info("Saving loans.Service state")
	snapshot := loans.Snapshot()
	if err := s.Write(&snapshot); err != nil {
		logrus.WithError(err).Error("error saving state")
	}
}
