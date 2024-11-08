package balancer

import (
	"work-balancer/internal/requester"
	"work-balancer/internal/worker"
)

type Pool []*worker.Worker

type Balancer interface {
	AddWorker(w *worker.Worker)
	Balance(work chan requester.Request)
	dispatch(req requester.Request)
	completed(w *worker.Worker)
}

func NewBalancer(done chan *worker.Worker) Balancer {
	return &balancer{}
}

type balancer struct {
	pool Pool
	done chan *worker.Worker
}

func (b *balancer) AddWorker(w *worker.Worker) {
	b.pool = append(b.pool, w)
}

func (b *balancer) Balance(work chan requester.Request) {
	for {
		select {
		case req := <-work:
			b.dispatch(req)
		case w := <-b.done:
			b.completed(w)
		}
	}
}

func (b *balancer) dispatch(req requester.Request) {
	w := b.pool[0]
	w.Requests <- req
	b.pool = append(b.pool[1:], w)
}

func (b *balancer) completed(w *worker.Worker) {
	b.pool = append(b.pool, w)
}
