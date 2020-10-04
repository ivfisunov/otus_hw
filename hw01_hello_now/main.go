package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	remoteHost := "ntp2.stratum2.ru"
	remotehostTime, err := ntp.Time(remoteHost)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	roundedRemotehostTime := remotehostTime.Round(time.Second)

	localhostTime := time.Now()
	roundedLocahostTime := localhostTime.Round(time.Second)

	fmt.Printf("current time: %v\n", roundedLocahostTime)
	fmt.Printf("exact time: %v\n", roundedRemotehostTime)
}
