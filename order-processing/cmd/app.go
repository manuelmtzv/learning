package main

import (
	"order-processing/internal/models"
)

func (app *application) run() {
	app.orderSimulator.Generate(app.ctx)
	pendingOrders := make(map[int]*models.Order)

	watchStream := app.watcher.Watch(app.ctx)
	pendingStream := app.manager.ManagePending(app.ctx, pendingOrders, watchStream)

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
