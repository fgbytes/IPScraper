package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
	fmt.Println("Main goroutine exit!")
	elapsed := time.Since(start)
	log.Printf("Checking  IP's took %s", elapsed)
}

func getIP(site string) string {

	start := time.Now()
	var ip string

	addrs, err := net.LookupIP(site)
	if err != nil {
		log.Print(err)
	}

	fmt.Println("looked up", len(addrs))
	if len(addrs) == 0 {
		ip = site + ",notfound"
	} else {
		ip = site + "," + addrs[0].String()
	}

	fmt.Println("received ip: ", ip)
	elapsed := time.Since(start)
	if elapsed > 1*time.Second {
		log.Printf("IP lookup took %s", elapsed)
	}
	return ip

}

func reader(raw chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(fileNameArgument)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}
	defer func() {
		errFileOpen := file.Close()
		if errFileOpen != nil {
			log.Fatal(errFileOpen)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		site := scanner.Text()
		raw <- site
	}

	// check for errors
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	isReadingFinished = true

}

func fixer(raw chan string, fixed chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	timeOutChan := make(chan string, 1)
	for {
		log.Println(len(raw), " ratio ", len(fixed))
		site := <-raw
		go func() { timeOutChan <- getIP(site) }()
		select {
		case receievedIP := <-timeOutChan:
			fixed <- receievedIP

			// use err and reply
		case <-time.After(100 * time.Millisecond):
			fixed <- site + ",timeout!"
			continue
		}

		if isReadingFinished && len(raw) == 0 {
			isTransferComplete = true
			break

		}
	}
}

func writer(fixed chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Create(fileNameArgument + "_result.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer func() {
		errFileOpen := file.Close()
		if errFileOpen != nil {
			log.Fatal(errFileOpen)
		}
	}()

	w := bufio.NewWriter(file)

	for {
		result := <-fixed
		fmt.Println(result)
		//	result += "/n"
		_, err := w.WriteString(result + "\n")
		fmt.Println("writer error: ", err)
		// time.Sleep(time.Second * 1)
		func() {
			err := w.Flush()
			if err != nil {
				log.Fatal("Error flushing file", err)
			}
		}()

		if isTransferComplete && len(fixed) == 0 {
			isTransferComplete = true
			break

		}

	}
}
