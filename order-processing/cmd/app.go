package main

import (
	"order-processing/internal/models"
	"sync"
	"time"
)

func (app *application) watch() (<-chan *models.Order, error) {
	pendingStream := make(chan *models.Order, 100)

	fetch := func() error {
		pending, err := app.store.Orders.GetCreatedOrders(app.ctx)
		if err != nil {
			app.logger.Errorf("Error fetching pending orders: %v", err)
			return err
		}

		app.logger.Infof("Orders watch query finished: %d new orders", len(pending))

		for _, order := range pending {
			select {
			case <-app.ctx.Done():
				return nil
			case pendingStream <- order:
			}
		}

		return nil
	}

	err := fetch()
	if err != nil {
		return pendingStream, err
	}

	go func() {
		defer close(pendingStream)
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-app.ctx.Done():
				return
			case <-ticker.C:
				fetch()
			}
		}
	}()

	return pendingStream, nil
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
				defer m.Unlock()
				if _, exists := pending[order.ID]; exists {
					continue
				}

				err := app.store.Orders.ChangeOrderStatus(app.ctx, order.ID, "pending")
				if err != nil {
					app.logger.Warnf("Error while setting order %d as pending: %v", order.ID, err)
					continue
				}

				pending[order.ID] = order
				pendingStream <- order
			}
		}
	}()

	return pendingStream
}

func (app *application) run() {
	pendingOrders := make(map[int]*models.Order)

	watchStream, err := app.watch()
	if err != nil {
		app.logger.Panic("Error on first query of new orders:", err)
	}

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
