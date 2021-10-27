package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func call(wg *sync.WaitGroup, domain string) {
	defer wg.Done()
	r, err := http.Get("http://"+domain)
	if err != nil {
		return
	}
	defer r.Body.Close()

	_, _ = ioutil.ReadAll(r.Body)

	fmt.Print(".")
}

func main() {
	var wg sync.WaitGroup

	args := os.Args
	if len(args)<3 {
		fmt.Println("Specify a domain and count as calling arguments\n\tExample: `goflood example.com 10`")
		return
	}

	domain := args[1]
	count, _ := strconv.Atoi(args[2])

	fmt.Println("Pinging "+domain+" "+strconv.Itoa(count)+" times")

	wg.Add(count)
	for i := 0; i < count; i++ {
		go call(&wg, domain)
	}
	wg.Wait()

	fmt.Println("Done")
}
