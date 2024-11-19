package main

type config struct {
	workers int
}

func main() {
	cfg := &config{
		workers: 10,
	}

	work := make(chan request)
	done := make(chan *worker)
	defer close(done)

	b := NewBalancer(done)

	for i := 0; i < cfg.workers; i++ {
		w := &worker{
			Requests: make(chan request),
			Index:    i,
		}
		b.AddWorker(w)

		go w.Work(done)
	}

	go b.Balance(work)

	req := NewRequester(cfg.workers)
	go req.Request(work)

	select {}
}
