package main

type Pool []*worker

type Balancer interface {
	AddWorker(w *worker)
	Balance(work chan request)
	dispatch(req request)
	completed(w *worker)
}

func NewBalancer(done chan *worker) Balancer {
	return &balancer{}
}

type balancer struct {
	pool Pool
	done chan *worker
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
	w := b.pool[0]
	w.Requests <- req
	b.pool = append(b.pool[1:], w)
}

func (b *balancer) completed(w *worker) {
	b.pool = append(b.pool, w)
}
