package main

import (
	"log"
	"time"
)

type Requester interface {
	Request(work chan<- request)
	logger(result int)
}

type request struct {
	Fn func() int
	C  chan int
}

func NewRequester(workers int) Requester {
	return &requester{
		workers: workers,
	}
}

type requester struct {
	workers int
}

func (r *requester) Request(work chan<- request) {
	c := make(chan int)

	for {
		time.Sleep(1 * time.Second)
		work <- request{workFn, c}
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
