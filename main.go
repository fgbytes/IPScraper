package main

import (
	"log"
	"os"
	"sync"
	"time"
)

var fileNameArgument string

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No input file specified. Shutting down")
	}
	fileNameArgument = os.Args[1]
	var raw = make(chan string, 128)
	var fixed = make(chan string, 128)
	var wgGlobal sync.WaitGroup
	wgGlobal.Add(2)
	var wgWorker sync.WaitGroup

	start := time.Now()

	go readerStart(raw, &wgGlobal)

	for w := 1; w <= 128; w++ {
		go worker(raw, fixed, &wgWorker, 3000)
		wgWorker.Add(1)
	}

	go writerStart(fixed, &wgGlobal)
	wgWorker.Wait()
	wgGlobal.Wait()

	elapsed := time.Since(start)
	log.Println(" Main goroutine exit!")
	log.Printf("Checking  IP's took %s", elapsed)
}
