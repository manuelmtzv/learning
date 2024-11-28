package main

import (
	"context"
	"order-processing/internal/store"
	"sync"

	"go.uber.org/zap"
)

func orderCleanup(store *store.Storage, logger *zap.SugaredLogger) {
	ctx := context.Background()

	logger.Info("Starting order cleanup")

	pending, err := store.Orders.GetOrdersByStatus(ctx, "pending")
	if err != nil {
		logger.Errorf("Error fetching pending orders: %v", err)
		return
	}

	var wg sync.WaitGroup

	for _, order := range pending {
		wg.Add(1)

		go func() {
			defer wg.Done()
			err := store.Orders.ChangeOrderStatus(ctx, order.ID, "created")
			if err != nil {
				logger.Warnf("Error while setting order %d as created: %v", order.ID, err)
			}
		}()
	}

	wg.Wait()

	logger.Info("Order cleanup finished")
}
