package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type config struct {
	randLimit     int
	delay         time.Duration
	tasksPerBlock int
}

func main() {
	cfg := &config{
		randLimit:     10000,
		delay:         2 * time.Second,
		tasksPerBlock: 5,
	}

	tasks := []int{}
	c := sync.Cond{
		L: &sync.Mutex{},
	}

	go addTasks(cfg, &tasks, &c)

	for {
		c.L.Lock()

		for len(tasks) == 0 {
			fmt.Println("Waiting for tasks...")
			c.Wait()
		}

		task := tasks[0]
		tasks = tasks[1:]
		fmt.Printf("Task %d have been processed\n", task)

		c.L.Unlock()
	}
}

func addTasks(cfg *config, tasks *[]int, c *sync.Cond) {
	for {
		time.Sleep(cfg.delay)
		newTasks := make([]int, cfg.tasksPerBlock)

		for i := 0; i < cfg.tasksPerBlock; i++ {
			randomNumber := rand.Intn(cfg.randLimit)
			newTasks[i] = randomNumber
			fmt.Println("Added task", randomNumber)
		}

		c.L.Lock()
		*tasks = append(*tasks, newTasks...)
		c.Signal()
		c.L.Unlock()
	}
}
