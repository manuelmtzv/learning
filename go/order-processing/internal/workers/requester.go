package workers

import (
	"context"
	"order-processing/internal/models"
	"order-processing/internal/store"

	"github.com/charmbracelet/log"
)

type Result struct {
	order   *models.Order
	success bool
}

type Request struct {
	order *models.Order
	c     chan *Result
}

type Requester interface {
	Request(context.Context, <-chan *models.Order, chan<- *Request, chan<- *models.Order)
}

type RequesterWorker struct {
	store  *store.Storage
	logger *log.Logger
}

func NewRequester(store *store.Storage, logger *log.Logger) Requester {
	return &RequesterWorker{
		store:  store,
		logger: logger,
	}
}

func (w *RequesterWorker) Request(ctx context.Context, pendingStream <-chan *models.Order, workStream chan<- *Request, processedStream chan<- *models.Order) {
	go func() {
		c := make(chan *Result)

		for {
			select {
			case <-ctx.Done():
				return
			case order := <-pendingStream:
				w.logger.Info("Pending order added:", "order", order)

				workStream <- &Request{order: order, c: c}

				go func() {
					select {
					case <-ctx.Done():
						return
					case result := <-c:
						if !result.success {
							w.logger.Errorf("Order %d could not be processed:", order.ID)
							return
						}
						processedStream <- result.order
					}
				}()
			}
		}
	}()
}
