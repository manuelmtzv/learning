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
	Generate()
}

type OrderSimulatorWorker struct {
	store  *store.Storage
	ctx    context.Context
	logger *zap.SugaredLogger
}

func NewOrderSimulator(ctx context.Context, store *store.Storage, logger *zap.SugaredLogger) OrderSimulator {
	return &OrderSimulatorWorker{
		store:  store,
		ctx:    ctx,
		logger: logger,
	}
}

func (w OrderSimulatorWorker) Generate() {
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
				w.store.Orders.CreateOrder(w.ctx, order)
			}
		}

		simulate()

		for {
			select {
			case <-w.ctx.Done():
				return
			case <-ticker.C:
				simulate()
			}
		}
	}()
}
