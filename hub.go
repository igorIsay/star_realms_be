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
}

func newHub() *Hub {
	return &Hub{
		action:     make(chan Action),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	log.Println("run state manager")

	deck := getDeck()
	stateManager := newStateManager(deck)
	middleware := newMiddleware(deck)

	go stateManager.run()
	go h.broadcast(&stateManager.json)
	pActions := middleware.preparareState()
	for _, a := range pActions {
		stateManager.action <- a
	}
	for {
		select {
		case client := <-h.register:
			log.Println("register")
			h.clients[client] = true
			stateManager.action <- &StateActionGetState{}
		case client := <-h.unregister:
			log.Println("unregister")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case action := <-h.action:
			actions := middleware.handle(string(action.message), action.client.playerId)
			for _, a := range actions {
				stateManager.action <- a
			}
		}
	}
}

func (h *Hub) broadcast(channel *chan []byte) {
	for {
		state := <-*channel
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
