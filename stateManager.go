package main

import (
	"encoding/json"
)

type StateManager struct {
	state  *State
	action chan StateAction
	json   chan []byte
}

type StateActionType int

const (
	ChangeCounterValue StateActionType = iota
	RandomCard
	MoveCard
	MoveAll
	GetState
	ChangeTurn
)

type PlayerId int

const (
	FirstPlayer PlayerId = iota + 1
	SecondPlayer
)

type Counter int

const (
	Trade Counter = iota
	Combat
	Authority
	Discard
)

type Operation int

const (
	Increase Operation = iota
	Decrease
	Set
)

type StateAction interface {
	Type() StateActionType
	Data() map[string]interface{}
}

type StateActionChangeCounterValue struct {
	player    PlayerId
	counter   Counter
	operation Operation
	value     int
}

func (s *StateActionChangeCounterValue) Type() StateActionType {
	return ChangeCounterValue
}

func (s *StateActionChangeCounterValue) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["player"] = s.player
	data["counter"] = s.counter
	data["operation"] = s.operation
	data["value"] = s.value
	return data
}

type StateActionRandomCard struct {
	from CardLocation
	to   CardLocation
}

func (s *StateActionRandomCard) Type() StateActionType {
	return RandomCard
}

func (s *StateActionRandomCard) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["from"] = s.from
	data["to"] = s.to
	return data
}

type StateActionMoveCard struct {
	id string
	to CardLocation
}

func (s *StateActionMoveCard) Type() StateActionType {
	return MoveCard
}

func (s *StateActionMoveCard) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = s.id
	data["to"] = s.to
	return data
}

type StateActionMoveAll struct {
	from CardLocation
	to   CardLocation
}

func (s *StateActionMoveAll) Type() StateActionType {
	return MoveAll
}

func (s *StateActionMoveAll) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["from"] = s.from
	data["to"] = s.to
	return data
}

type StateActionGetState struct{}

func (s *StateActionGetState) Type() StateActionType {
	return GetState
}

func (s *StateActionGetState) Data() map[string]interface{} {
	data := make(map[string]interface{})
	return data
}

type StateActionChangeTurn struct{}

func (s *StateActionChangeTurn) Type() StateActionType {
	return ChangeTurn
}

func (s *StateActionChangeTurn) Data() map[string]interface{} {
	data := make(map[string]interface{})
	return data
}

func newStateManager() *StateManager {
	return &StateManager{
		state:  newState(),
		action: make(chan StateAction),
		json:   make(chan []byte),
	}
}

func (s *StateManager) run() {
	for {
		action := <-s.action
		switch action.Type() {
		case ChangeCounterValue:
			data := action.Data()
			player := data["player"].(PlayerId)
			counter := data["counter"].(Counter)
			operation := data["operation"].(Operation)
			value := data["value"].(int)
			c := &s.state.FirstPlayerCounters
			if player == SecondPlayer {
				c = &s.state.SecondPlayerCounters
			}
			counters := make(map[Counter]*int)
			counters[Trade] = &c.Trade
			counters[Authority] = &c.Authority
			counters[Combat] = &c.Combat
			counters[Discard] = &c.Discard
			calc(counters[counter], value, operation)
		case RandomCard:
			data := action.Data()
			from := data["from"].(CardLocation)
			to := data["to"].(CardLocation)
			deck := s.cardsByLocation(from)
			if len(deck) == 0 {
				if from == FirstPlayerDeck {
					discard := s.cardsByLocation(FirstPlayerDiscard)
					for _, c := range discard {
						c.Location = FirstPlayerDeck
					}
				}
				if from == SecondPlayerDeck {
					discard := s.cardsByLocation(SecondPlayerDiscard)
					for _, c := range discard {
						c.Location = SecondPlayerDeck
					}
				}
				deck = s.cardsByLocation(from)
			}
			card, ok := randomCard(deck)
			if ok {
				card.Location = to
			} else {
				//TODO: handle exception
			}
		case MoveCard:
			data := action.Data()
			id := data["id"].(string)
			to := data["to"].(CardLocation)
			card, ok := s.cardById(id)
			if ok {
				card.Location = to
			}
		case MoveAll:
			data := action.Data()
			from := data["from"].(CardLocation)
			to := data["to"].(CardLocation)
			cards := s.cardsByLocation(from)
			for _, card := range cards {
				card.Location = to
			}
		case ChangeTurn:
			if s.state.Turn == FirstPlayer {
				s.state.Turn = SecondPlayer
			} else {
				s.state.Turn = FirstPlayer
			}
		case GetState:
			state, _ := json.Marshal(s.state)
			s.json <- state
		}
	}
}

func (s *StateManager) cardById(id string) (*Card, bool) {
	card, ok := s.state.Cards[id]
	return card, ok
}

func (s *StateManager) cardsByLocation(l CardLocation) map[string]*Card {
	result := make(map[string]*Card)
	for id, card := range s.state.Cards {
		if card.Location == l {
			result[id] = card
		}
	}
	return result
}

func randomCard(deck map[string]*Card) (*Card, bool) {
	for _, card := range deck {
		return card, true
	}
	return nil, false
}

func calc(a *int, b int, operation Operation) {
	switch operation {
	case Increase:
		*a += b
	case Decrease:
		*a -= b
	case Set:
		*a = b
	}
}
