package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const EOF = "...EOF"

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

	eofCh := make(chan string, 1)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	go func() {
		<-ch
		os.Exit(-1)
	}()

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "...connected to %v\n", net.JoinHostPort(host, port))

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go receive(client, wg, eofCh)
	go send(client, wg, eofCh)
	wg.Wait()
}

func receive(client TelnetClient, wg *sync.WaitGroup, eofCh chan string) {
	defer wg.Done()
	err := client.Receive()
	if err != nil {
		select {
		case <-eofCh:
			return
		default:
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}

func send(client TelnetClient, wg *sync.WaitGroup, eofCh chan string) {
	defer wg.Done()
	defer client.Close()
	err := client.Send()
	if err != nil {
		fmt.Fprint(os.Stderr, "...connection closed by peer\n")
		return
	}
	fmt.Fprintln(os.Stderr, EOF)
	eofCh <- EOF
}
