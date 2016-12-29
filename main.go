package main

import (
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		log.Fatal("No input file specified. Shutting down")
	}
	//raw data channel
	var raw = make(chan string, chanSise)
	//fixed - channel for sites with checked IP
	var fixed = make(chan string, chanSise)
	//WaitGroup for reader and writer
	var wgGlobal sync.WaitGroup
	//WaitGroup wor workers in the pool
	var wgWorker sync.WaitGroup
	wgGlobal.Add(2)

	// starting reader, starting reader pool with W as worker count
	go readerStart(raw, &wgGlobal)
	for w := 1; w <= workerCount; w++ {
		go worker(raw, fixed, &wgWorker, maxDelay)
		wgWorker.Add(1)
	}
	go writerStart(fixed, &wgGlobal)

	wgWorker.Wait()
	close(fixed)
	wgGlobal.Wait()

	elapsed := time.Since(start)
	log.Println(" Main goroutine exit!")
	log.Printf("Checking  IP's took %s", elapsed)
}
