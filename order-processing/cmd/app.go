package main

import (
	"math/rand"
	"order-processing/internal/models"
	"time"
)

func (app *application) orderSimulate() {
	go func() {
		ticker := time.NewTicker(time.Duration(rand.Intn(5)+2) * time.Second)
		defer ticker.Stop()

		simulate := func() {
			amount := rand.Intn(300) + 1

			app.logger.Infof("Adding %v new simulated orders", amount)
			for i := 0; i <= amount; i++ {
				order := &models.Order{
					Status: "created",
				}
				app.store.Orders.CreateOrder(app.ctx, order)
			}
		}

		simulate()

		for {
			select {
			case <-app.ctx.Done():
				return
			case <-ticker.C:
				simulate()
			}
		}
	}()
}

func (app *application) run() {
	app.orderSimulate()
	pendingOrders := make(map[int]*models.Order)

	watchStream := app.watcher.Watch()
	pendingStream := app.manager.ManagePending(pendingOrders, watchStream)

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
