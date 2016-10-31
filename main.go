package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
	//	"time"
)

var IsReadingFinished bool = false
var IsTransferComplete bool = false
var IsWritingComplete bool = false
var fileNameArgument = os.Args[1]

func main() {
	start := time.Now()

	var raw = make(chan string, 256)
	var fixed = make(chan string, 300)
	var wg sync.WaitGroup

	wg.Add(3)

	go reader(raw, &wg)
	go fixer(raw, fixed, &wg)
	go writer(fixed, &wg)

	// var input string
	// fmt.Scanln(&input)
	wg.Wait()
	fmt.Println("Main goroutine exit!")
	elapsed := time.Since(start)
	log.Printf("Checking  IP's took %s", elapsed)
}

func getIP(site string) string {

	start := time.Now()
	var ip string

	addrs, _ := net.LookupIP(site)

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
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		site := scanner.Text()
		raw <- site
	}

	// check for errors
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	IsReadingFinished = true

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

		if IsReadingFinished && len(raw) == 0 {
			IsTransferComplete = true
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
	defer file.Close()

	w := bufio.NewWriter(file)

	for {

		result := <-fixed
		fmt.Println(result)
		//	result += "/n"
		_, err := w.WriteString(result + "\n")
		fmt.Println("writer error: ", err)
		// time.Sleep(time.Second * 1)
		w.Flush()

		if IsTransferComplete && len(fixed) == 0 {
			IsTransferComplete = true
			break

		}

	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
