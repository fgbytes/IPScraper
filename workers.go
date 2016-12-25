package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

func readerStart(raw chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileNameArgument)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			log.Fatal(errClose)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		site := scanner.Text()
		raw <- site
	}

	// check for errors
	// if err = scanner.Err(); err != nil {
	// 	log.Fatal(err)
	//}
	close(raw)
	log.Println("reader finished execution")
}

func writerStart(fixed chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Create(fileNameArgument + "_result.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("Error closing file", err)
		}
	}()

	w := bufio.NewWriter(file)

	for result := range fixed {
		log.Printf("filewriter is writing%s", result)

		_, err := w.WriteString(result + "\n")
		if err != nil {
			log.Printf("writer error: %s", err)
		}
		func() {
			if err := w.Flush(); err != nil {
				log.Fatal("Error flushing file", err)
			}
		}()

	}

}

func worker(raw chan string, fixed chan string, wg *sync.WaitGroup, delay int) {
	defer wg.Done()
	c := make(chan string, 1)
	for site := range raw {
		go func() { c <- getIP(site) }()
		select {
		case ip := <-c:
			fixed <- site + "," + ip
		case <-time.After(time.Duration(delay) * time.Millisecond):
			// call timed out
			fixed <- site + "," + "timedout!"
		}

	}
	close(fixed)
}
