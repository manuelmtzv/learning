package worker

import "work-balancer/internal/requester"

type Worker struct {
	Requests chan requester.Request
	Pending  int
	Index    int
}

func (w *Worker) Work(done chan *Worker) {
	for {
		req := <-w.Requests
		req.C <- req.Fn()
		done <- w
	}
}
