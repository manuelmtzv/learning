package workers

import (
	"context"
	"order-processing/internal/models"
	"order-processing/internal/store"

	"go.uber.org/zap"
)

type Request struct {
	order *models.Order
	c     chan *models.Order
}

type Requester interface {
	Request(context.Context, <-chan *models.Order, chan<- *Request, chan<- *models.Order)
}

type RequesterWorker struct {
	store  *store.Storage
	logger *zap.SugaredLogger
}

func NewRequester(store *store.Storage, logger *zap.SugaredLogger) Requester {
	return &RequesterWorker{
		store:  store,
		logger: logger,
	}
}

func (w *RequesterWorker) Request(ctx context.Context, pendingStream <-chan *models.Order, workStream chan<- *Request, processedStream chan<- *models.Order) {
	go func() {
		c := make(chan *models.Order, 100)
		defer close(c)

		for {
			select {
			case <-ctx.Done():
				w.logger.Info("Shutting down requester worker gracefully.")
				return
			case order := <-pendingStream:
				w.logger.Info("Pending order added:", order)

				workStream <- &Request{order: order, c: c}

				select {
				case result := <-c:
					processedStream <- result
				case <-ctx.Done():
					w.logger.Info("Context done, exiting worker.")
					return
				}
			}
		}
	}()
}
