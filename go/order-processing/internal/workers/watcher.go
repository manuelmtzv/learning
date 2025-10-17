package workers

import (
	"context"
	"order-processing/internal/models"
	"order-processing/internal/store"
	"time"

	"github.com/charmbracelet/log"
)

type Watcher interface {
	Watch(ctx context.Context) <-chan *models.Order
}

type WatcherWorker struct {
	store  *store.Storage
	logger *log.Logger
}

func NewWatcher(store *store.Storage, logger *log.Logger) Watcher {
	return &WatcherWorker{
		store:  store,
		logger: logger,
	}
}

func (w *WatcherWorker) Watch(ctx context.Context) <-chan *models.Order {
	pendingStream := make(chan *models.Order, 2000)

	fetchPendingOrders := func() {
		pending, err := w.store.Orders.GetOrdersByStatus(ctx, "created")
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

	go fetchPendingOrders()

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
