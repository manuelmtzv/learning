package main

import (
	"order-processing/internal/models"
	"order-processing/internal/workers"
)

func (app *application) run(simulate bool) {
	if simulate {
		app.orderSimulator.Generate(app.ctx)
	}
	pendingOrders := make(map[int]*models.Order)

	watchStream := app.watcher.Watch(app.ctx)
	pendingStream := app.manager.ManagePending(app.ctx, pendingOrders, watchStream)

	workStream := make(chan *workers.Request, 2000)
	processedStream := make(chan *models.Order, 1000)

	app.requester.Request(app.ctx, pendingStream, workStream, processedStream)

	app.balancer.Balance(app.ctx, workStream)

	go func() {
		for {
			select {
			case <-app.ctx.Done():
				return
			case order := <-processedStream:
				app.logger.Info("Order processed!:", "order", order)
			}
		}
	}()
}
