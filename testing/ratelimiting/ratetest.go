package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:8080/healthz"
	requests := 40                     
	interval := 10 * time.Millisecond

	for i := 1; i <= requests; i++ {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("[%2d] Error: %v\n", i, err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		fmt.Printf("[%2d] Status: %d, Body: %s\n", i, resp.StatusCode, string(body))
		time.Sleep(interval)
	}
}
