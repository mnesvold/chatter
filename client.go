package main

import (
	"io"
	"log"
)

type client struct {
	rwc  io.ReadWriteCloser
	send <-chan []byte
	recv chan<- []byte
}

func newClient(rwc io.ReadWriteCloser, send <-chan []byte, recv chan<- []byte) *client {
	client := client{rwc, send, recv}
	go client.read()
	return &client
}

func (c *client) read() {
	buf := make([]byte, 256)
	for {
		n, err := c.rwc.Read(buf)
		if err == io.EOF {
      break
    } else if err != nil {
			log.Fatal(err) // TODO: handle error
		}
		_ = n
		c.recv <- buf[:n]
	}
}
