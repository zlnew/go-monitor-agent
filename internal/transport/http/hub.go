// Package http
package http

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"zlnew/monitor-agent/internal/core"
	"zlnew/monitor-agent/internal/infra/logger"
)

type Hub struct {
	rooms map[string]*Room
	mu    sync.Mutex
	log   logger.Logger
}

type Room struct {
	name       string
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
	log        logger.Logger
}

type ClientMessage struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

func NewHub(log logger.Logger) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		log:   log,
	}
}

func (h *Hub) GetOrCreateRoom(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, exists := h.rooms[name]; exists {
		return room
	}

	room := &Room{
		name:       name,
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 5),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		log:        h.log,
	}
	h.rooms[name] = room

	go room.Run()
	return room
}

func (h *Hub) BroadcastMetrics(metrics core.Metrics) {
	room := h.GetOrCreateRoom("metrics")
	msg := map[string]interface{}{
		"channel": "metrics",
		"payload": metrics,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		h.log.Error("marshal metrics", "error", err)
		return
	}
	room.broadcast <- bytes
}

func (r *Room) Run() {
	r.log.Info("room started", "room", r.name)
	for {
		select {
		case client := <-r.register:
			r.mu.Lock()
			r.clients[client] = true
			r.mu.Unlock()
			r.log.Info("client registered", "room", r.name, "remote_addr", client.RemoteAddr())
		case client := <-r.unregister:
			r.mu.Lock()
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				client.Close()
				r.log.Info("client unregistered", "room", r.name, "remote_addr", client.RemoteAddr())
			}
			r.mu.Unlock()
		case message := <-r.broadcast:
			r.mu.Lock()
			for client := range r.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					r.log.Error("write message", "error", err)
					go func(c *websocket.Conn) {
						r.unregister <- c
					}(client)
				}
			}
			r.mu.Unlock()
		}
	}
}
