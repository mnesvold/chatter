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
	connections   []*websocket.Conn
	broadcastChan chan []byte
}

func NewServer() (server *Server) {
	server = &Server{broadcastChan: make(chan []byte)}
	go server.broadcast()
	return
}

func (s *Server) broadcast() {
	for {
		payload := <-s.broadcastChan
		for _, conn := range s.connections {
			conn.Write(payload)
		}
	}
}

func (s *Server) HandleConnection(ws *websocket.Conn) {
	s.connections = append(s.connections, ws)
	decoder := json.NewDecoder(ws)
	for {
		var payload map[string]interface{}
		if err := decoder.Decode(&payload); err == io.EOF {
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

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	websocket.Handler(s.HandleConnection).ServeHTTP(w, req)
}
