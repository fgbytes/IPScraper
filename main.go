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

	var raw = make(chan string, 256)
	var fixed = make(chan string, 300)
	var wg sync.WaitGroup

	wg.Add(3)

	go reader(raw, &wg)
	go fixer(raw, fixed, &wg)
	go writer(fixed, &wg)

	wg.Wait()
	log.Println("Main goroutine exit!")
	elapsed := time.Since(start)
	log.Printf("Checking  IP's took %s", elapsed)
}
