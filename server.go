package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
)

const (
	protocol = "chatter.clocktower.systems"
)

type Server struct {
	connections     []io.ReadWriteCloser
	broadcastChan   chan []byte
	newClientChan   chan io.ReadWriteCloser
	closeClientChan chan io.ReadWriteCloser
}

func NewServer() (server *Server) {
	server = &Server{
		broadcastChan:   make(chan []byte),
		newClientChan:   make(chan io.ReadWriteCloser),
		closeClientChan: make(chan io.ReadWriteCloser),
	}
	go server.handleClients()
	return
}

func (s *Server) handleClients() {
	for {
		select {
		case payload := <-s.broadcastChan:
			for _, conn := range s.connections {
				go conn.Write(payload)
			}
		case client := <-s.newClientChan:
			s.connections = append(s.connections, client)
			go s.readClient(client)
		case client := <-s.closeClientChan:
			err := client.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (s *Server) readClient(ws io.ReadWriteCloser) {
	decoder := json.NewDecoder(ws)
	for {
		var payload map[string]interface{}
		if err := decoder.Decode(&payload); err == io.EOF {
			s.closeClientChan <- ws
			break
		} else if err != nil {
			log.Fatal(err)
		}
		message := payload["message"]

		response := make(map[string]interface{})
		response["message"] = message
		responsePayload, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		s.broadcastChan <- responsePayload
	}
}

func (s *Server) HandleConnection(ws *websocket.Conn) {
	s.newClientChan <- ws
	select {} // returning from this method would close the connection, so block instead
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	websocket.Handler(s.HandleConnection).ServeHTTP(w, req)
}
