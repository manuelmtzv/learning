package http

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var client = &http.Client{Timeout: 5 * time.Second}

func CheckStatus(wg *sync.WaitGroup, url string) {
	defer wg.Done()
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("%s\t %s %s\n", time.Now().Format(time.Stamp), url, err.Error())
		return
	}
	defer resp.Body.Close()

	fmt.Printf("%s\t %s %s\n", time.Now().Format(time.Stamp), url, resp.Status)
}
