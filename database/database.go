// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package database

import (
	"context"
	"time"

	"github.com/elgopher/deebee-loans/service"
	"github.com/elgopher/deebee/store"
)

func StartLoans(ctx context.Context, s Store) (loans *SynchronizedLoans, done <-chan struct{}, err error) {
	loans, err = loadState(ctx, s)
	if err != nil {
		return nil, nil, err
	}

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		for {
			select {
			case <-time.After(time.Minute):
				saveState(ctx, loans, s)
			case <-ctx.Done():
				saveState(ctx, loans, s)
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

func loadState(ctx context.Context, s Store) (*SynchronizedLoans, error) {
	snapshot := service.Snapshot{}
	version, err := s.ReadLatest(&snapshot)
	if store.IsVersionNotFound(err) {
		log.WithError(ctx, err).Warn("No snapshot found")
	} else if err != nil {
		return nil, err
	} else {
		log.With(ctx, "version", version).Info("Snapshot loaded")
	}

	loans := service.FromSnapshot(snapshot)

	return &SynchronizedLoans{
		loans: loans,
	}, nil
}

func saveState(ctx context.Context, loans *SynchronizedLoans, s Store) {
	log.Info(ctx, "Saving loans.Service state")
	snapshot := loans.Snapshot()
	if err := s.Write(&snapshot); err != nil {
		log.WithError(ctx, err).Error("error saving state")
	}
}
