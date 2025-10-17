package main

import (
	"fmt"
	"sync"
)

type Button struct {
	Clicked *sync.Cond
}

func main() {
	subscribe := func(stream <-chan interface{}, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			goroutineRunning.Done()
			<-stream
			fn()
		}()

		goroutineRunning.Wait()
	}

	buttonClicked := make(chan interface{})
	wg := sync.WaitGroup{}

	wg.Add(3)

	subscribe(buttonClicked, func() {
		fmt.Println("Maximizing window.")
		wg.Done()
	})
	subscribe(buttonClicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		wg.Done()
	})
	subscribe(buttonClicked, func() {
		fmt.Println("Mouse clicked.")
		wg.Done()
	})

	close(buttonClicked)
	wg.Wait()
}

// func main() {
// 	button := Button{
// 		Clicked: sync.NewCond(&sync.Mutex{}),
// 	}

// 	subscribe := func(c *sync.Cond, fn func()) {
// 		var goroutineRunning sync.WaitGroup
// 		goroutineRunning.Add(1)

// 		go func() {
// 			goroutineRunning.Done()
// 			c.L.Lock()
// 			defer c.L.Unlock()
// 			c.Wait()
// 			fn()
// 		}()

// 		goroutineRunning.Wait()
// 	}

// 	var clickRegistered sync.WaitGroup
// 	clickRegistered.Add(3)

// 	subscribe(button.Clicked, func() {
// 		fmt.Println("Maximizing window.")
// 		clickRegistered.Done()
// 	})

// 	subscribe(button.Clicked, func() {
// 		fmt.Println("Displaying annoying dialog box!")
// 		clickRegistered.Done()
// 	})

// 	subscribe(button.Clicked, func() {
// 		fmt.Println("Mouse clicked.")
// 		clickRegistered.Done()
// 	})

// 	button.Clicked.Broadcast()
// 	clickRegistered.Wait()
// }
