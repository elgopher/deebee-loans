package etcd

import (
	"context"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func WaitToBeElected(ctx context.Context) (*Leader, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	session, err := concurrency.NewSession(etcdClient)
	if err != nil {
		return nil, err
	}

	election := concurrency.NewElection(session, "/leader-election/")
	if err = election.Campaign(ctx, "e"); err != nil {
		return nil, err
	}
	return &Leader{
		etcdClient: etcdClient,
		session:    session,
		election:   election,
	}, nil
}

type Leader struct {
	etcdClient *clientv3.Client
	session    *concurrency.Session
	election   *concurrency.Election
}

func (l *Leader) Resign() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l.election.Resign(ctx)
	l.session.Close()
	l.etcdClient.Close()
}
