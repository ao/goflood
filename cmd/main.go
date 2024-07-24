package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorOrange = "\033[33m"
)

func main() {
	fmt.Println(Banner())

	var targetURL string
	var concurrency int
	var duration time.Duration
	var showHelp bool

	flag.StringVar(&targetURL, "url", "", "Target URL")
	flag.IntVar(&concurrency, "n", 5, "Concurrency Step")
	flag.DurationVar(&duration, "t", 10*time.Second, "Duration for which the program should run")
	flag.BoolVar(&showHelp, "help", false, "Show help")

	flag.Parse()

	if targetURL == "" {
		fmt.Println("Please specify a target URL")
		flag.Usage()
		os.Exit(0)
	}

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	var wg sync.WaitGroup
	var counter int32
	var totalRequests int32
	var successfulRequests int32
	var failedRequests int32
	var requestTimes []time.Duration
	var mu sync.Mutex

	// Interrupt handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		wg.Wait() // Wait for ongoing requests to finish
		os.Exit(0)
	}()

	startTime := time.Now()

	timer := time.NewTimer(duration)
	defer timer.Stop()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	fmt.Printf("Running for %s\n\n", duration)

	go func() {
		for range ticker.C {
			elapsedTime := time.Since(startTime)
			remainingTime := duration - elapsedTime
			fmt.Printf("\nTime elapsed: %ds - Time remaining: %ds\n", int(elapsedTime.Seconds()), int(remainingTime.Seconds()))

			intervalWg := sync.WaitGroup{}

			for i := 0; i < concurrency; i++ {
				wg.Add(1)
				intervalWg.Add(1)
				atomic.AddInt32(&counter, 1)
				atomic.AddInt32(&totalRequests, 1)
				go func() {
					makeRequest(targetURL, &counter, &wg, &successfulRequests, &failedRequests, &requestTimes, &mu)
					intervalWg.Done()
				}()
			}
			intervalWg.Wait() // Wait for all requests in this interval to finish
		}
	}()

	<-timer.C // Wait for the duration to elapse

	// Stop the ticker
	ticker.Stop()

	// Wait until all requests are done
	wg.Wait()

	elapsed := time.Since(startTime)

	// Print summary report in tabular format
	printSummaryReport(targetURL, concurrency, elapsed, totalRequests, successfulRequests, failedRequests, requestTimes)
}

func printSummaryReport(targetURL string, concurrency int, elapsed time.Duration, totalRequests, successfulRequests, failedRequests int32, requestTimes []time.Duration) {
	durationStr := elapsed.Round(time.Millisecond).String()

	// Calculate min, median, and max request times
	sort.Slice(requestTimes, func(i, j int) bool {
		return requestTimes[i] < requestTimes[j]
	})
	minTime := requestTimes[0]
	maxTime := requestTimes[len(requestTimes)-1]
	medianTime := requestTimes[len(requestTimes)/2]

	// Determine the maximum length of the target URL and set the width accordingly
	titlesWidth := 25
	maxWidth := 40
	if len(targetURL) > maxWidth {
		maxWidth = len(targetURL)
	}

	// Create dynamic separator
	separator := createSeparator(titlesWidth, maxWidth)

	fmt.Println("\n*** Summary Report ***")
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*s |\n", titlesWidth, "Target URL", maxWidth, targetURL)
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*d |\n", titlesWidth, "Concurrency", maxWidth, concurrency)
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*s |\n", titlesWidth, "Duration", maxWidth, durationStr)
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*d |\n", titlesWidth, "Total Requests", maxWidth, totalRequests)
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*d |\n", titlesWidth, "Successful Requests", maxWidth, successfulRequests)
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*d |\n", titlesWidth, "Failed Requests", maxWidth, failedRequests)
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*s |\n", titlesWidth, "Min Request Time", maxWidth, minTime)
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*s |\n", titlesWidth, "Median Request Time", maxWidth, medianTime)
	fmt.Println(separator)
	fmt.Printf("| %-*s | %-*s |\n", titlesWidth, "Max Request Time", maxWidth, maxTime)
	fmt.Println(separator)
}

func createSeparator(labelWidth, valueWidth int) string {
	return fmt.Sprintf("+-%s-+-%s-+", dashes(labelWidth), dashes(valueWidth))
}

func dashes(n int) string {
	ret := ""
	for i := 0; i < n; i++ {
		ret += "-"
	}
	return ret
}
