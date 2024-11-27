package workers

import (
	"context"
	"order-processing/internal/models"
	"order-processing/internal/store"
	"sync"

	"go.uber.org/zap"
)

type Manager interface {
	ManagePending(map[int]*models.Order, <-chan *models.Order) <-chan *models.Order
}

type ManagerWorker struct {
	store  *store.Storage
	ctx    context.Context
	logger *zap.SugaredLogger
}

func NewManager(ctx context.Context, store *store.Storage, logger *zap.SugaredLogger) Manager {
	return &ManagerWorker{
		store:  store,
		ctx:    ctx,
		logger: logger,
	}
}

func (w *ManagerWorker) ManagePending(pending map[int]*models.Order, watchStream <-chan *models.Order) <-chan *models.Order {
	pendingStream := make(chan *models.Order)

	go func() {
		m := &sync.Mutex{}

		for {
			select {
			case <-w.ctx.Done():
				return
			case order := <-watchStream:
				m.Lock()
				if _, exists := pending[order.ID]; exists {
					m.Unlock()
					continue
				}

				err := w.store.Orders.ChangeOrderStatus(w.ctx, order.ID, "pending")
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
