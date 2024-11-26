package main

import (
	"order-processing/internal/db"
	"order-processing/internal/env"
	"order-processing/internal/store"
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
			workers: env.GetInt("WORKERS", 4),
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

	app := &application{
		processor: cfg.processor,
		store:     store,
		logger:    logger,
	}

	app.run()
}
