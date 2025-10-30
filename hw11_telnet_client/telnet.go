package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrNotConnected = errors.New("no connection, call Connect() first")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type BaseTelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	connection                   net.Conn
	clientScanner, serverScanner *bufio.Scanner
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &BaseTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *BaseTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.connection = conn
	c.clientScanner = bufio.NewScanner(c.in)
	c.serverScanner = bufio.NewScanner(c.connection)

	return nil
}

func (c *BaseTelnetClient) Close() error {
	if c.connection == nil {
		return ErrNotConnected
	}

	err := c.connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *BaseTelnetClient) Send() error {
	if c.connection == nil {
		return ErrNotConnected
	}

	if !c.clientScanner.Scan() {
		return io.EOF
	}

	_, err := c.connection.Write([]byte(fmt.Sprintf("%s\n", c.clientScanner.Bytes())))
	if err != nil {
		return err
	}

	return nil
}

func (c *BaseTelnetClient) Receive() error {
	if c.connection == nil {
		return ErrNotConnected
	}

	if !c.serverScanner.Scan() {
		return io.EOF
	}

	_, err := fmt.Fprintln(c.out, c.serverScanner.Text())
	if err != nil {
		return err
	}

	return nil
}
