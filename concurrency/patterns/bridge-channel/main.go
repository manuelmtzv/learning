package main

import "fmt"

func orDone(done, c <-chan any) <-chan any {
	valueStream := make(chan any)
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

func bridge(done <-chan any, chanStream <-chan <-chan any) <-chan any {
	valueStream := make(chan any)
	go func() {
		defer close(valueStream)
		for {
			var stream <-chan any
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) {
				select {
				case valueStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valueStream
}

func main() {
	getPagedResults := func() <-chan <-chan any {
		chanStream := make(chan (<-chan any))
		go func() {
			defer close(chanStream)
			for page := 1; page <= 3; page++ {
				pageChan := make(chan any, 2)
				pageChan <- fmt.Sprintf("Page %d - A", page)
				pageChan <- fmt.Sprintf("Page %d - B", page)
				close(pageChan)
				chanStream <- pageChan
			}
		}()
		return chanStream
	}

	done := make(chan any)
	defer close(done)

	for result := range bridge(done, getPagedResults()) {
		fmt.Println(result)
	}

}
