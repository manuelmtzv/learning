package main

type worker struct {
	Requests chan request
	Pending  int
	Index    int
}

func (w *worker) Work(done chan *worker) {
	for {
		req := <-w.Requests
		req.C <- req.Fn()
		done <- w
	}
}
