package main

import (
	"context"
	"order-processing/internal/workers"
	"time"

	"go.uber.org/zap"
)

type application struct {
	logger         *zap.SugaredLogger
	ctx            context.Context
	watcher        workers.Watcher
	manager        workers.Manager
	orderSimulator workers.OrderSimulator
}

type config struct {
	processor processorConfig
	db        dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type processorConfig struct {
	workers int
}
