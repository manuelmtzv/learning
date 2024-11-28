package main

import (
	"context"
	"order-processing/internal/db"
	"order-processing/internal/env"
	"order-processing/internal/store"
	"order-processing/internal/workers"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()

	err := env.Load()
	if err != nil {
		logger.Panic(err)
	}

	cfg := &config{
		processor: processorConfig{
			workers:        env.GetInt("WORKERS", 4),
			simulateOrders: env.GetBool("SIMULATE_ORDERS", false),
		},
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Panic(err)
	}

	defer db.Close()
	logger.Infow("DB connected")

	store := store.NewStorage(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orderSimulator := workers.NewOrderSimulator(store, logger)
	watcher := workers.NewWatcher(store, logger)
	manager := workers.NewManager(store, logger)
	requester := workers.NewRequester(store, logger)
	balancer := workers.NewBalancer(logger, store, cfg.processor.workers)

	app := &application{
		logger:         logger,
		ctx:            ctx,
		watcher:        watcher,
		manager:        manager,
		orderSimulator: orderSimulator,
		requester:      requester,
		balancer:       balancer,
	}

	signalStream := make(chan os.Signal, 1)
	signal.Notify(signalStream, os.Interrupt, syscall.SIGTERM)

	go app.run(cfg.processor.simulateOrders)

	<-signalStream
	logger.Infof("Shutting down gracefully... CTRL + C to force.")

	cancel()

	orderCleanup(store, logger)

	time.Sleep(1 * time.Second) // Simulate cleanup
	logger.Infof("Application stopped.")
}
