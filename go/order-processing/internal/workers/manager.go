package workers

import (
	"context"
	"order-processing/internal/models"
	"order-processing/internal/store"
	"sync"

	"github.com/charmbracelet/log"
)

type Manager interface {
	ManagePending(context.Context, map[int]*models.Order, <-chan *models.Order) <-chan *models.Order
}

type ManagerWorker struct {
	store  *store.Storage
	logger *log.Logger
}

func NewManager(store *store.Storage, logger *log.Logger) Manager {
	return &ManagerWorker{
		store:  store,
		logger: logger,
	}
}

func (w *ManagerWorker) ManagePending(ctx context.Context, pending map[int]*models.Order, watchStream <-chan *models.Order) <-chan *models.Order {
	pendingStream := make(chan *models.Order, 2000)

	go func() {
		defer close(pendingStream)
		m := &sync.Mutex{}

		for {
			select {
			case <-ctx.Done():
				return
			case order := <-watchStream:
				m.Lock()
				if _, exists := pending[order.ID]; exists {
					m.Unlock()
					continue
				}

				err := w.store.Orders.ChangeOrderStatus(ctx, order.ID, "pending")
				if err != nil {
					m.Unlock()
					w.logger.Warnf("Error while setting order %d as pending: %v", order.ID, err)
					continue
				}

				pending[order.ID] = order
				m.Unlock()

				pendingStream <- order
			}
		}
	}()

	return pendingStream
}
