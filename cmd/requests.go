package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func makeRequest(url string, counter *int32, wg *sync.WaitGroup, successfulRequests *int32, failedRequests *int32, requestTimes *[]time.Duration, mu *sync.Mutex) {
	defer wg.Done()
	defer atomic.AddInt32(counter, -1)

	start := time.Now()
	res, err := http.Get(url)
	elapsed := time.Since(start)

	mu.Lock()
	*requestTimes = append(*requestTimes, elapsed)
	mu.Unlock()

	if err != nil {
		fmt.Printf("%sError making request: %v (Duration: %s)%s\n", colorOrange, err, elapsed, colorReset)
		atomic.AddInt32(failedRequests, 1)
	} else {
		fmt.Printf("%sRequest successful (Duration: %s) %d%s\n", colorGreen, elapsed, res.StatusCode, colorReset)
		atomic.AddInt32(successfulRequests, 1)
		res.Body.Close()
	}
}
