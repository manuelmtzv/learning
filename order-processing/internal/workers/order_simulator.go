package workers

import (
	"context"
	"math/rand"
	"order-processing/internal/models"
	"order-processing/internal/store"
	"time"

	"go.uber.org/zap"
)

type OrderSimulator interface {
	Generate(ctx context.Context)
}

type OrderSimulatorWorker struct {
	store  *store.Storage
	logger *zap.SugaredLogger
}

func NewOrderSimulator(store *store.Storage, logger *zap.SugaredLogger) OrderSimulator {
	return &OrderSimulatorWorker{
		store:  store,
		logger: logger,
	}
}

func (w OrderSimulatorWorker) Generate(ctx context.Context) {
	simulate := func() {
		amount := rand.Intn(800) + 1

		w.logger.Infof("Adding %v new simulated orders", amount)
		for i := 0; i <= amount; i++ {
			order := &models.Order{
				Status: "created",
			}
			w.store.Orders.CreateOrder(ctx, order)
		}
	}

	go func() {
		ticker := time.NewTicker(time.Duration(rand.Intn(80)+30) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				simulate()
			}
		}
	}()
}
