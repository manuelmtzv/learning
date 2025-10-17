package main

import (
	"context"
	"order-processing/internal/workers"
	"time"

	"github.com/charmbracelet/log"
)

type application struct {
	logger *log.Logger
	ctx    context.Context

	watcher        workers.Watcher
	manager        workers.Manager
	orderSimulator workers.OrderSimulator
	requester      workers.Requester
	balancer       workers.Balancer
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
	workers        int
	simulateOrders bool
}
