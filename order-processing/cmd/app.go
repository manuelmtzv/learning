package main

import (
	"order-processing/internal/models"
	"time"
)

func (app *application) fetchPendingOrders() <-chan *models.Order {
	pendingStream := make(chan *models.Order, 100)

	go func() {
		defer close(pendingStream)
		ticker := time.NewTicker(5 * time.Second)
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
			}
		}
	}()

	return pendingStream
}

func (app *application) run() {
	pendingStream := app.fetchPendingOrders()

	go func() {
		for {
			select {
			case <-app.ctx.Done():
				return
			case order := <-pendingStream:
				app.logger.Info("Processing order:", order)
			}
		}
	}()
}
