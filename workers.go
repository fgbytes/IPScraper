package main

import (
	"bufio"
	"log"
	"os"
	"sync"
)

func readFileJob(raw chan string, wg *sync.WaitGroup) {
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
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	isReadingFinished = true
	close(raw)

}

func worker(raw <-chan string, fixed chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for site := range raw {
		timeOutChan := make(chan string, 1)
		go func() { timeOutChan <- getIP(site) }()

		select {
		case receievedIP := <-timeOutChan:
			fixed <- receievedIP

		}

		if isReadingFinished && len(raw) == 0 {
			isTransferComplete = true
			break

		}
	}
}

// func checkIPjob(raw chan string, fixed chan string, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	timeOutChan := make(chan string, 1)
// 	for {
// 		site := <-raw
// 		go func() { timeOutChan <- getIP(site) }()

// 		select {
// 		case receievedIP := <-timeOutChan:
// 			fixed <- receievedIP
// 		case <-time.After(100 * time.Millisecond):
// 			fixed <- site + ",timeout!"
// 			continue
// 		}

// 		if isReadingFinished && len(raw) == 0 {
// 			isTransferComplete = true
// 			break

// 		}
// 	}
// }

func writeFileJob(fixed chan string, wg *sync.WaitGroup) {
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

	for {
		result := <-fixed
		log.Println(result)
		//	result += "/n"
		_, err := w.WriteString(result + "\n")
		if err != nil {
			log.Printf("writer error: %s", err)
		}
		func() {
			if err := w.Flush(); err != nil {
				log.Fatal("Error flushing file", err)
			}
		}()

		if isTransferComplete && len(fixed) == 0 {
			isTransferComplete = true
			break

		}

	}
}
