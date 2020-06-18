package main

import (
	"log"
)

// Hub maintains the set of active clients and process actions
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound action from the clients.
	action chan Action

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	stateManager *StateManager
}

func newHub() *Hub {
	return &Hub{
		action:       make(chan Action),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
		stateManager: newStateManager(),
	}
}

func (h *Hub) run() {
	log.Println("run state manager")
	go h.stateManager.run()
	go h.broadcast()
	middleware := newMiddleware()
	pActions := middleware.preparareState()
	for _, a := range pActions {
		h.stateManager.action <- a
	}
	for {
		select {
		case client := <-h.register:
			log.Println("register")
			h.clients[client] = true
			h.stateManager.action <- &StateActionGetState{}
		case client := <-h.unregister:
			log.Println("unregister")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case action := <-h.action:
			actions := middleware.handle(string(action.message), action.client.playerId)
			for _, a := range actions {
				h.stateManager.action <- a
			}
		}
	}
}

func (h *Hub) broadcast() {
	for {
		state := <-h.stateManager.json
		for client := range h.clients {
			select {
			case client.send <- state:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}
