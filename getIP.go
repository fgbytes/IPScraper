package main

import (
	"log"
	"net"
	"time"
)

func getIP(site string) string {

	start := time.Now()
	var namedIPResponce string

	siteIP, err := net.LookupIP(site)
	if err != nil {
		log.Print(err)
	}

	if len(siteIP) == 0 {
		namedIPResponce = site + ",notfound"
	} else {
		namedIPResponce = site + "," + siteIP[0].String()
	}
	log.Println("received ip: ", namedIPResponce)

	if elapsed := time.Since(start); elapsed > 1*time.Second {
		log.Printf("IP lookup took %s", elapsed)
	}
	return namedIPResponce

}
