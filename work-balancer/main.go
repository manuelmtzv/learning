package main

import (
	"work-balancer/internal/balancer"
	"work-balancer/internal/requester"
	"work-balancer/internal/worker"
)

type config struct {
	workers int
}

func main() {
	cfg := &config{
		workers: 10,
	}

	work := make(chan requester.Request)
	done := make(chan *worker.Worker)

	b := balancer.NewBalancer(done)

	for i := 0; i < cfg.workers; i++ {
		w := &worker.Worker{
			Requests: make(chan requester.Request),
			Index:    i,
		}
		b.AddWorker(w)

		go w.Work(done)
	}

	go b.Balance(work)

	req := requester.NewRequester(cfg.workers)
	go req.Request(work)

	select {}
}
