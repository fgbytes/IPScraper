package main

import (
	"log"
	"net"
	"time"
)

func getIP(site string) string {

	start := time.Now()
	var ip string

	addrs, err := net.LookupIP(site)
	if err != nil {
		log.Print(err)
	}

	if len(addrs) == 0 {
		ip = site + ",notfound"
	} else {
		ip = site + "," + addrs[0].String()
	}

	log.Println("received ip: ", ip)
	elapsed := time.Since(start)
	if elapsed > 1*time.Second {
		log.Printf("IP lookup took %s", elapsed)
	}
	return ip

}
