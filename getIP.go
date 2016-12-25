package main

import (
	"log"
	"net"
	"time"
)

func getIP(site string) string {

	start := time.Now()
	IP, err := net.LookupIP(site)
	if err != nil {
		log.Println(err)
	}

	if len(IP) == 0 {
		return "not_found"
	}

	log.Printf("received ip : %s", IP[0].String())

	if elapsed := time.Since(start); elapsed > 10*time.Millisecond {
		log.Printf("IP lookup took %s", elapsed)
	}

	return IP[0].String()
}
