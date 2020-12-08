package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeoutFlag = flag.String("timeout", "10s", "connection timeout")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("provide host and port")
	}
	host := args[0]
	port := args[1]

	timeout, err := time.ParseDuration(*timeoutFlag)
	if err != nil {
		log.Fatal("error parsing time")
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	go func() {
		<-ch
		os.Exit(-1)
	}()

	client := NewTelnetClient(host+":"+port, timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "...connected to %v\n", host+":"+port)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go Receive(client, wg)
	go Send(client, wg)
	wg.Wait()
}

func Receive(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()
	err := client.Receive()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return
	}
}

func Send(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()
	err := client.Send()
	if err != nil {
		fmt.Fprint(os.Stderr, "...connection closed by peer\n")
		return
	}
	fmt.Fprint(os.Stderr, "...EOF\n")
	client.Close()
}
