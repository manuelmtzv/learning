package workers

import (
	"context"
	"order-processing/internal/models"
	"order-processing/internal/store"
	"time"

	"go.uber.org/zap"
)

type Watcher interface {
	Watch() <-chan *models.Order
}

type WatcherWorker struct {
	store  *store.Storage
	ctx    context.Context
	logger *zap.SugaredLogger
}

func NewWatcher(ctx context.Context, store *store.Storage, logger *zap.SugaredLogger) Watcher {
	return &WatcherWorker{
		store:  store,
		ctx:    ctx,
		logger: logger,
	}
}

func (w WatcherWorker) Watch() <-chan *models.Order {
	pendingStream := make(chan *models.Order, 500)

	fetch := func() {
		pending, err := w.store.Orders.GetCreatedOrders(w.ctx)
		if err != nil {
			w.logger.Errorf("Error fetching pending orders: %v", err)
			return
		}

		w.logger.Infof("Orders watch query finished: %d new orders", len(pending))

		for _, order := range pending {
			select {
			case <-w.ctx.Done():
				return
			case pendingStream <- order:
			}
		}
	}

	fetch()

	go func() {
		defer close(pendingStream)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-w.ctx.Done():
				return
			case <-ticker.C:
				fetch()
			}
		}
	}()

	return pendingStream
}
