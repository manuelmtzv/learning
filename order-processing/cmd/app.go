package main

import (
	"order-processing/internal/models"
	"sync"
	"time"
)

func (app *application) fetchPendingOrders() <-chan *models.Order {
	pendingStream := make(chan *models.Order, 100)

	go func() {
		defer close(pendingStream)
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-app.ctx.Done():
				return
			case <-ticker.C:
				pending, err := app.store.Orders.GetPendingOrders(app.ctx)
				if err != nil {
					app.logger.Errorf("Error fetching pending orders: %v", err)
					continue
				}

				if len(pending) == 0 {
					app.logger.Info("No orders available")
				}

				for _, order := range pending {
					select {
					case <-app.ctx.Done():
						return
					case pendingStream <- order:
					}
				}

				app.logger.Info("Orders refresh finished")
			}
		}
	}()

	return pendingStream
}

func (app *application) pendingOrdersStream(pending map[int]*models.Order, fetchStream <-chan *models.Order) <-chan *models.Order {
	pendingStream := make(chan *models.Order)

	go func() {
		m := &sync.Mutex{}

		for {
			select {
			case <-app.ctx.Done():
				return
			case order := <-fetchStream:
				m.Lock()
				if _, exists := pending[order.ID]; exists {
					m.Unlock()
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

func (app *application) run() {
	pendingOrders := make(map[int]*models.Order)

	fetchStream := app.fetchPendingOrders()
	pendingStream := app.pendingOrdersStream(pendingOrders, fetchStream)

	go func() {
		for {
			select {
			case <-app.ctx.Done():
				return
			case order := <-pendingStream:
				app.logger.Info("Pending order added:", order)
			}
		}
	}()
}
