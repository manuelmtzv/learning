package requester

import (
	"log"
	"time"
)

type Request struct {
	Fn func() int
	C  chan int
}

type Requester interface {
	Request(work chan<- Request)
	logger(result int)
}

func NewRequester(workers int) Requester {
	return &requester{
		workers: workers,
	}
}

type requester struct {
	workers int
}

func (r *requester) Request(work chan<- Request) {
	c := make(chan int)

	for {
		time.Sleep(1 * time.Second)
		work <- Request{workFn, c}
		result := <-c
		r.logger(result)
	}
}

func (r *requester) logger(result int) {
	log.Printf("result: %d\n", result)
}

func workFn() int {
	return 1
}
