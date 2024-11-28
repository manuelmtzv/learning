package workers

import (
	"context"
	"order-processing/internal/models"
	"order-processing/internal/store"

	"go.uber.org/zap"
)

type Processor interface {
	Work(context.Context)
}

type ProcessorWorker struct {
	logger   *zap.SugaredLogger
	store    *store.Storage
	requests chan *Request
	pending  int
	index    int
}

func NewProcessor(logger *zap.SugaredLogger, store *store.Storage) Processor {
	return &ProcessorWorker{
		logger:   logger,
		store:    store,
		requests: make(chan *Request),
	}
}

func (w *ProcessorWorker) Work(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case req := <-w.requests:
				order, err := w.process(ctx, req.order)
				success := true

				if err != nil {
					w.logger.Error("Failed to process order:", err)
					success = false
				}

				req.c <- &Result{
					order:   order,
					success: success,
				}
			}
		}
	}()
}

func (w *ProcessorWorker) process(ctx context.Context, order *models.Order) (*models.Order, error) {
	w.store.Orders.ChangeOrderStatus(ctx, order.ID, "processed")
	return order, nil
}
