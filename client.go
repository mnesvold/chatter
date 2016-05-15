package main

import (
	"encoding/json"
	"io"
	"log"
)

type client struct {
	rwc  io.ReadWriteCloser
	send <-chan []byte
	recv chan<- map[string]interface{}
}

func newClient(rwc io.ReadWriteCloser, send <-chan []byte, recv chan<- map[string]interface{}) *client {
	client := client{rwc, send, recv}
	go client.read()
	go client.write()
	return &client
}

func (c *client) read() {
	decoder := json.NewDecoder(c.rwc)
	for {
		var data map[string]interface{}
		err := decoder.Decode(&data)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err) // TODO: handle error
		}
		c.recv <- data
	}
}

func (c *client) write() {
	for {
		message := <-c.send
		_, err := c.rwc.Write(message)
		if err != nil {
			log.Fatal(err) // TODO: handle error
		}
	}
}
