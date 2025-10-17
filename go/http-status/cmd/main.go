package main

import (
	"errors"
	"http-status/internal/http"
	"http-status/internal/url"
	"log"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal(errors.New("missing urls"))
	}

	urls, err := url.ExplodeUrls(args[0])
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)

		go http.CheckStatus(&wg, url)
	}

	wg.Wait()
}
