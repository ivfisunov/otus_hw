package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ErrEOF = errors.New("...EOF")

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

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT)

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "...connected to %v\n", net.JoinHostPort(host, port))

	ctx, cancel := context.WithCancel(context.Background())

	defer client.Close()
	go receive(client, cancel)
	go send(client, cancel)

	select {
	case <-stopCh:
		return
	case <-ctx.Done():
		return
	}
}

func receive(client TelnetClient, cancel context.CancelFunc) {
	err := client.Receive()
	if err != nil {
		cancel()
		return
	}
}

func send(client TelnetClient, cancel context.CancelFunc) {
	err := client.Send()
	if err != nil {
		fmt.Fprint(os.Stderr, "...connection closed by peer\n")
		cancel()
		return
	}
	fmt.Fprintln(os.Stderr, ErrEOF)
	cancel()
}
