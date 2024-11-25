package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	ch2 := make(chan int)
	go func() {
		for {
			select {
			case val := <-ch:
				fmt.Println("Received:", val)
			case val := <-ch2:
				fmt.Println("Received:", val)
			}
		}
	}()

	for i := 0; i < 5; i++ {
		ch <- i
		time.Sleep(1 * time.Second)
	}

	for i := 5; i < 10; i++ {
		ch2 <- i
		time.Sleep(1 * time.Second)
	}
}
