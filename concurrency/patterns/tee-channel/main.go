package main

import "fmt"

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valueStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valueStream
}

func tee(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func main() {
	done := make(chan interface{})
	in := make(chan interface{})
	defer close(done)

	go func() {
		defer close(in)
		for i := 0; i < 5; i++ {
			in <- i
		}
	}()

	out1, out2 := tee(done, in)

	for val := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val, <-out2)
	}
}
