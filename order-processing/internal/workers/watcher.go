package workers

import (
	"context"
	"order-processing/internal/models"
	"order-processing/internal/store"
	"time"

	"go.uber.org/zap"
)

type Watcher interface {
	Watch(ctx context.Context) <-chan *models.Order
}

type WatcherWorker struct {
	store  *store.Storage
	logger *zap.SugaredLogger
}

func NewWatcher(store *store.Storage, logger *zap.SugaredLogger) Watcher {
	return &WatcherWorker{
		store:  store,
		logger: logger,
	}
}

func (w WatcherWorker) Watch(ctx context.Context) <-chan *models.Order {
	pendingStream := make(chan *models.Order, 500)

	fetchPendingOrders := func() {
		pending, err := w.store.Orders.GetCreatedOrders(ctx)
		if err != nil {
			w.logger.Errorf("Error fetching pending orders: %v", err)
			return
		}

		w.logger.Infof("Orders watch query finished: %d new orders", len(pending))

		for _, order := range pending {
			select {
			case <-ctx.Done():
				return
			case pendingStream <- order:
			}
		}
	}

	fetchPendingOrders()

	go func() {
		defer close(pendingStream)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fetchPendingOrders()
			}
		}
	}()

	return pendingStream
}
