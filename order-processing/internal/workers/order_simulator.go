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
	go func() {
		ticker := time.NewTicker(time.Duration(rand.Intn(5)+2) * time.Second)
		defer ticker.Stop()

		simulate := func() {
			amount := rand.Intn(300) + 1

			w.logger.Infof("Adding %v new simulated orders", amount)
			for i := 0; i <= amount; i++ {
				order := &models.Order{
					Status: "created",
				}
				w.store.Orders.CreateOrder(ctx, order)
			}
		}

		simulate()

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
