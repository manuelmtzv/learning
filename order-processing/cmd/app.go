package main

import (
	"order-processing/internal/models"
	"sync"
	"time"
)

func (app *application) watch() <-chan *models.Order {
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
				pending, err := app.store.Orders.GetCreatedOrders(app.ctx)
				if err != nil {
					app.logger.Errorf("Error fetching pending orders: %v", err)
					continue
				}

				if len(pending) == 0 {
					app.logger.Info("No orders available to enter in process")
				}

				for _, order := range pending {
					select {
					case <-app.ctx.Done():
						return
					case pendingStream <- order:
					}
				}

				app.logger.Info("Orders watch process finished")
			}
		}
	}()

	return pendingStream
}

func (app *application) managePending(pending map[int]*models.Order, watchStream <-chan *models.Order) <-chan *models.Order {
	pendingStream := make(chan *models.Order)

	go func() {
		m := &sync.Mutex{}

		for {
			select {
			case <-app.ctx.Done():
				return
			case order := <-watchStream:
				m.Lock()
				if _, exists := pending[order.ID]; exists {
					m.Unlock()
					continue
				}

				err := app.store.Orders.ChangeOrderStatus(app.ctx, order.ID, "pending")
				if err != nil {
					app.logger.Error("Error while setting order as pending:", err)
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

	watchStream := app.watch()
	pendingStream := app.managePending(pendingOrders, watchStream)

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
