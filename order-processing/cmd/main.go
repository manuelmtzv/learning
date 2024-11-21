package main

type config struct {
	orderWorkers int
}

func main() {
	cfg := &config{
		orderWorkers: 10,
	}

}
