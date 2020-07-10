package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

type StateManager struct {
	state  *State
	action chan StateAction
	json   chan []byte
}

type StateActionType int

const (
	ChangeCounterValue StateActionType = iota
	TopCard
	MoveCard
	MoveAll
	GetState
	ChangeTurn
	RequestUserAction
	AddActivatedAbility
	DisableActivatedAbility
	ResetActivatedAbilities
	ShuffleDeck
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
	ShipsOnTop
	fleetFlag
	blobs
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

type StateActionTopCard struct {
	from CardLocation
	to   CardLocation
}

func (s *StateActionTopCard) Type() StateActionType {
	return TopCard
}

func (s *StateActionTopCard) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["from"] = s.from
	data["to"] = s.to
	return data
}

type StateActionShuffleDeck struct {
	deck CardLocation
}

func (s *StateActionShuffleDeck) Type() StateActionType {
	return ShuffleDeck
}

func (s *StateActionShuffleDeck) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["deck"] = s.deck
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

type StateActionRequestUserAction struct {
	player PlayerId
	action UserAction
	cardId string
}

func (s *StateActionRequestUserAction) Type() StateActionType {
	return RequestUserAction
}

func (s *StateActionRequestUserAction) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["player"] = s.player
	data["action"] = s.action
	data["cardId"] = s.cardId
	return data
}

type StateActionAddActivatedAbility struct {
	cardId    string
	abilityId AbilityId
}

func (s *StateActionAddActivatedAbility) Type() StateActionType {
	return AddActivatedAbility
}

func (s *StateActionAddActivatedAbility) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["cardId"] = s.cardId
	data["abilityId"] = s.abilityId
	return data
}

type StateActionDisableActivatedAbility struct {
	cardId    string
	abilityId AbilityId
}

func (s *StateActionDisableActivatedAbility) Type() StateActionType {
	return DisableActivatedAbility
}

func (s *StateActionDisableActivatedAbility) Data() map[string]interface{} {
	data := make(map[string]interface{})
	data["cardId"] = s.cardId
	data["abilityId"] = s.abilityId
	return data
}

func newStateManager(deck *map[string]*CardEntry) *StateManager {
	return &StateManager{
		state:  newState(deck),
		action: make(chan StateAction),
		json:   make(chan []byte),
	}
}

type StateActionResetActivatedAbilities struct{}

func (s *StateActionResetActivatedAbilities) Type() StateActionType {
	return ResetActivatedAbilities
}

func (s *StateActionResetActivatedAbilities) Data() map[string]interface{} {
	data := make(map[string]interface{})
	return data
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
			counters[ShipsOnTop] = &c.ShipsOnTop
			counters[fleetFlag] = &c.fleetFlag
			counters[blobs] = &c.blobs
			calc(counters[counter], value, operation)
		case TopCard:
			data := action.Data()
			from := data["from"].(CardLocation)
			to := data["to"].(CardLocation)
			deck := s.cardsByLocation(from)
			if len(deck) == 0 && (from == FirstPlayerDeck || from == SecondPlayerDeck) {
				if from == FirstPlayerDeck {
					discard := s.cardsByLocation(FirstPlayerDiscard)
					for _, c := range discard {
						s.state.lastIndex[c.Location] -= 1
						s.state.lastIndex[FirstPlayerDeck] += 1
						c.Location = FirstPlayerDeck
						c.index = s.state.lastIndex[FirstPlayerDeck]
					}
				}
				if from == SecondPlayerDeck {
					discard := s.cardsByLocation(SecondPlayerDiscard)
					for _, c := range discard {
						s.state.lastIndex[c.Location] -= 1
						s.state.lastIndex[SecondPlayerDeck] += 1
						c.Location = SecondPlayerDeck
						c.index = s.state.lastIndex[SecondPlayerDeck]
					}
				}
				deck = s.cardsByLocation(from)
				s.shuffle(deck)
			}
			for _, card := range deck {
				if card.index == s.state.lastIndex[card.Location] {
					s.state.lastIndex[card.Location] -= 1
					s.state.lastIndex[to] += 1
					card.Location = to
					card.index = s.state.lastIndex[to]
					break
				}
			}
		case ShuffleDeck:
			data := action.Data()
			deck := data["deck"].(CardLocation)
			s.shuffle(s.cardsByLocation(deck))
		case MoveCard:
			data := action.Data()
			id := data["id"].(string)
			to := data["to"].(CardLocation)
			card, ok := s.cardById(id)
			if ok {
				s.state.lastIndex[card.Location] -= 1
				s.state.lastIndex[to] += 1
				card.Location = to
				card.index = s.state.lastIndex[to]
			}
		case MoveAll:
			data := action.Data()
			from := data["from"].(CardLocation)
			to := data["to"].(CardLocation)
			cards := s.cardsByLocation(from)
			for _, card := range cards {
				s.state.lastIndex[card.Location] -= 1
				s.state.lastIndex[to] += 1
				card.Location = to
				card.index = s.state.lastIndex[to]
			}
		case ChangeTurn:
			if s.state.Turn == FirstPlayer {
				s.state.Turn = SecondPlayer
			} else {
				s.state.Turn = FirstPlayer
			}
		case RequestUserAction:
			data := action.Data()
			player := data["player"].(PlayerId)
			userAction := data["action"].(UserAction)
			cardId := data["cardId"].(string)
			actionRequest := ActionRequest{
				Action: userAction,
				CardId: cardId,
			}
			switch player {
			case FirstPlayer:
				s.state.FirstPlayerActionRequest = actionRequest
			case SecondPlayer:
				s.state.SecondPlayerActionRequest = actionRequest
			}
		case AddActivatedAbility:
			data := action.Data()
			cardId := data["cardId"].(string)
			abilityId := data["abilityId"].(AbilityId)
			abilities, ok := s.state.ActivatedAbilities[cardId]
			if !ok {
				s.state.ActivatedAbilities[cardId] = make(map[AbilityId]bool)
				abilities, _ = s.state.ActivatedAbilities[cardId]
			}
			abilities[abilityId] = true
		case DisableActivatedAbility:
			data := action.Data()
			cardId := data["cardId"].(string)
			abilityId := data["abilityId"].(AbilityId)
			abilities, ok := s.state.ActivatedAbilities[cardId]
			if ok {
				abilities[abilityId] = false
			}
		case ResetActivatedAbilities:
			s.state.ActivatedAbilities = make(map[string]ActivatedAbilities)
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

func (s *StateManager) shuffle(deck map[string]*Card) {
	indexes := []int{}
	for _, card := range deck {
		indexes = append(indexes, card.index)
	}
	rand.Seed(time.Now().Unix())
	for _, card := range deck {
		idx := rand.Intn(len(indexes))
		card.index = indexes[idx]
		indexes = removeFromSlice(indexes, idx)
	}

}

func removeFromSlice(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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
