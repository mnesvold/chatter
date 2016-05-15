package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

const (
	protocol = "chatter.clocktower.systems"
)

type Server struct {
	clients         []*client
	recvChan        chan map[string]interface{}
	broadcastChan   chan []byte
	newClientChan   chan *client
	closeClientChan chan *client
}

func NewServer() (server *Server) {
	server = &Server{
		recvChan:        make(chan map[string]interface{}),
		broadcastChan:   make(chan []byte),
		newClientChan:   make(chan *client),
		closeClientChan: make(chan *client),
	}
	go server.handleClients()
	return
}

func (s *Server) handleClients() {
	for {
		select {
		case payload := <-s.recvChan:
			go s.receive(payload)
		case payload := <-s.broadcastChan:
			for _, client := range s.clients {
				go client.Send(payload)
			}
		case client := <-s.newClientChan:
			s.clients = append(s.clients, client)
		case client := <-s.closeClientChan:
			err := client.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (s *Server) receive(payload map[string]interface{}) {
	message := payload["message"]

	reply := map[string]interface{}{"message": message}
	replyPayload, err := json.Marshal(reply)
	if err != nil {
		log.Printf("marshalling %q failed: %v", reply, err)
	} else {
		s.broadcastChan <- replyPayload
	}
}

func (s *Server) HandleConnection(ws *websocket.Conn) {
	s.newClientChan <- newClient(ws, make(chan []byte), s.recvChan)
	select {} // returning from this method would close the connection, so block instead
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	websocket.Handler(s.HandleConnection).ServeHTTP(w, req)
}
