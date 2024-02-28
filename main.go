package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func Banner() string {
	return `
  ▄████  ▒█████    █████▒██▓     ▒█████   ▒█████  ▓█████▄ 
 ██▒ ▀█▒▒██▒  ██▒▓██   ▒▓██▒    ▒██▒  ██▒▒██▒  ██▒▒██▀ ██▌
▒██░▄▄▄░▒██░  ██▒▒████ ░▒██░    ▒██░  ██▒▒██░  ██▒░██   █▌
░▓█  ██▓▒██   ██░░▓█▒  ░▒██░    ▒██   ██░▒██   ██░░▓█▄   ▌
░▒▓███▀▒░ ████▓▒░░▒█░   ░██████▒░ ████▓▒░░ ████▓▒░░▒████▓ 
 ░▒   ▒ ░ ▒░▒░▒░  ▒ ░   ░ ▒░▓  ░░ ▒░▒░▒░ ░ ▒░▒░▒░  ▒▒▓  ▒ 
  ░   ░   ░ ▒ ▒░  ░     ░ ░ ▒  ░  ░ ▒ ▒░   ░ ▒ ▒░  ░ ▒  ▒ 
░ ░   ░ ░ ░ ░ ▒   ░ ░     ░ ░   ░ ░ ░ ▒  ░ ░ ░ ▒   ░ ░  ░ 
      ░     ░ ░             ░  ░    ░ ░      ░ ░     ░    
                                                   ░
`
}

func makeRequest(targetURL string, counter *int32, wg *sync.WaitGroup) {
	defer atomic.AddInt32(counter, -1)
	defer wg.Done()

	start := time.Now()

	resp, err := http.Get(targetURL)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Process response if needed

	if resp.StatusCode == http.StatusOK {
		// Successful response handling
		// You can add more metrics or processing here
	}

	// Print request duration
	fmt.Printf("Request to %s took %s\n", targetURL, time.Since(start))
}

func main() {
	fmt.Println(Banner())

	var targetURL string
	var baseConcurrency, concurrencyStep int
	var duration time.Duration
	var showHelp bool

	flag.StringVar(&targetURL, "url", "", "Target URL")
	flag.IntVar(&baseConcurrency, "concurrency", 10, "Base Concurrency")
	flag.IntVar(&concurrencyStep, "step", 1, "Concurrency Step")
	flag.DurationVar(&duration, "duration", 10*time.Second, "Duration for which the program should run")
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

	for concurrency := baseConcurrency; concurrency <= baseConcurrency+concurrencyStep; concurrency += concurrencyStep {
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			atomic.AddInt32(&counter, 1)
			go makeRequest(targetURL, &counter, &wg)
		}
	}

	<-timer.C // Wait for the specified duration

	// Wait until all requests are done
	for atomic.LoadInt32(&counter) > 0 {
		time.Sleep(100 * time.Millisecond)
	}

	elapsed := time.Since(startTime)

	fmt.Printf("Target URL: %s\n", targetURL)
	fmt.Printf("Base Concurrency: %d\n", baseConcurrency)
	fmt.Printf("Concurrency Step: %d\n", concurrencyStep)
	fmt.Printf("Duration: %s\n", elapsed)
}
