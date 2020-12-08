package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type client struct {
	connection net.Conn
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{address: address, timeout: timeout, in: in, out: out}
}

func (c *client) Connect() error {
	var err error
	c.connection, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Close() error {
	return c.connection.Close()
}

func (c *client) Send() error {
	_, err := io.Copy(c.connection, c.in)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.connection)
	if err != nil {
		return err
	}
	return nil
}
