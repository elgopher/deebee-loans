// (c) 2021 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jacekolszak/deebee-loans/database"
	"github.com/jacekolszak/deebee-loans/service"
	"github.com/jacekolszak/deebee-loans/web"
	"github.com/jacekolszak/deebee/codec"
	"github.com/jacekolszak/deebee/json"
	"github.com/jacekolszak/deebee/replicator"
	"github.com/jacekolszak/deebee/store"
	"github.com/jacekolszak/yala/adapter/printer"
	"github.com/jacekolszak/yala/logger"
)

func main() {
	logger.SetAdapter(printer.StdoutAdapter()) // configure logging

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	mainDir := flag.String("mainDir", "/tmp/loans", "Directory where data will be stored")
	backupDir := flag.String("backupDir", "/tmp/loans-backup", "Directory where data will be replicated once per hour")
	flag.Parse()
	logger.With(ctx, "main_data_dir", *mainDir).With("backup_data_dir", *backupDir).Info("Opening store")

	s := openStore(ctx, *mainDir, *backupDir)

	loans, done, err := database.StartLoans(ctx, s)
	if err != nil {
		logger.WithError(ctx, err).Error("Starting loans failed")
		os.Exit(1)
	}
	defer func() {
		<-done // wait until database saves snapshot of service.Loans
	}()

	if err := web.ListenAndServe(ctx, loans); err != nil && err != http.ErrServerClosed {
		logger.WithError(ctx, err).Error("Error starting server")
	}
}

func openStore(ctx context.Context, mainDir, backupDir string) *replicatedJsonStore {
	mainStore, err := store.Open(mainDir)
	if err != nil {
		logger.WithError(ctx, err).Error("error opening DeeBee store")
		os.Exit(1)
	}
	backupStore, err := store.Open(backupDir)
	if err != nil {
		logger.WithError(ctx, err).Error("error opening DeeBee backup store")
		os.Exit(1)
	}

	go func() {
		err = replicator.StartFromTo(ctx, mainStore, backupStore, replicator.Interval(time.Hour))
		if err != nil {
			logger.WithError(ctx, err).Error("cannot start replication") // continue even though replication does not work
		}
	}()

	return &replicatedJsonStore{
		mainStore:   mainStore,
		backupStore: backupStore,
	}
}

type replicatedJsonStore struct {
	mainStore   *store.Store
	backupStore *store.Store
}

func (a *replicatedJsonStore) ReadLatest(out *service.Snapshot) (store.Version, error) {
	return replicator.ReadLatest(json.Decoder(out), a.mainStore, a.backupStore)
}

func (a *replicatedJsonStore) Write(in *service.Snapshot, options ...store.WriterOption) error {
	return codec.Write(a.mainStore, json.Encoder(in), options...)
}
