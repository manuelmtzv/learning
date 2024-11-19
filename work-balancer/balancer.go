package main

import "fmt"

type Pool []*worker

type Balancer interface {
	AddWorker(w *worker)
	Balance(work chan request)
	dispatch(req request)
	completed(w *worker)
}

type balancer struct {
	pool Pool
	done chan *worker
}

func NewBalancer(done chan *worker) Balancer {
	return &balancer{
		pool: make(Pool, 0),
		done: done,
	}
}

func (b *balancer) AddWorker(w *worker) {
	b.pool = append(b.pool, w)
}

func (b *balancer) Balance(work chan request) {
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *balancer) dispatch(req request) {
	if len(b.pool) == 0 {
		fmt.Println("No workers available to dispatch")
		return
	}
	w := b.pool[0]
	b.pool = b.pool[1:]
	w.Requests <- req
}

func (b *balancer) completed(w *worker) {
	b.pool = append(b.pool, w)
}
