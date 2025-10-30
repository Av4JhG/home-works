package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var host, port string
	var timeout time.Duration
	logger := log.New(os.Stderr, "telnet:", 0)

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connect timeout")
	flag.Parse()
	host = flag.Arg(0)
	port = flag.Arg(1)

	sigintChan := make(chan os.Signal, 1)
	signal.Notify(sigintChan, syscall.SIGINT)

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	connectErr := client.Connect()
	if connectErr != nil {
		logger.Fatal(connectErr)
	}
	defer client.Close()

	receiveChan := make(chan error)
	go func(doneCh chan<- error) {
		for {
			receiveErr := client.Receive()
			if receiveErr != nil {
				close(doneCh)
				logger.Println(receiveErr)

				return
			}
		}
	}(receiveChan)

	sendChan := make(chan error)
	go func(doneCh chan<- error) {
		for {
			sendErr := client.Send()
			if sendErr != nil {
				close(doneCh)
				logger.Println(sendErr)

				return
			}
		}
	}(sendChan)

	for {
		select {
		case <-receiveChan:
			logger.Println("Connection was closed by server")
			return
		case <-sendChan:
			logger.Println("Connection was closed by peer")

			return
		case <-sigintChan:
			logger.Println("SIGINT. Closing.")

			return
		}
	}
}
