package main

import (
	"log"
	"os"
	"sync"
	"time"
)

var isReadingFinished = false
var isTransferComplete = false
var fileNameArgument string

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No input file specified. Shutting down")
	}
	fileNameArgument = os.Args[1]
	start := time.Now()

	var raw = make(chan string, 8)
	var fixed = make(chan string, 8)
	var wg sync.WaitGroup

	wg.Add(3 + 1024)

	go readFileJob(raw, &wg)
	//go checkIPjob(raw, fixed, &wg)
	for w := 1; w <= 1024; w++ {
		go worker(raw, fixed, &wg)
	}

	go writeFileJob(fixed, &wg)

	wg.Wait()
	elapsed := time.Since(start)
	log.Println(" Main goroutine exit!")
	log.Printf("Checking  IP's took %s", elapsed)
}
