package main

import (
	"strconv"
	"strings"
)

type Middleware struct {
	deck      map[string]*CardEntry
	allyState *AllyState
}

type AllyState struct {
	flags     map[Faction]bool
	abilities map[Faction]Abilities
}

const (
	TradeRowQty             int = 5
	HandCardsQty            int = 5
	FirstPlayerHandCardsQty int = 3
)

type Faction int

const (
	Unaligned Faction = iota
	Blob
	MachineCult
	StarEmpire
	TradeFederation
)

type Abilities []*Ability

type CardEntry struct {
	cost                 int
	faction              Faction
	primaryAbilities     Abilities
	utilizationAbilities Abilities
	allyAbilities        Abilities
}

type Ability struct {
	player PlayerPointer
	action func(PlayerId) StateAction
}

type PlayerPointer int

const (
	Current PlayerPointer = iota
	Opponent
)

var emptyUtilizationAbilities []*Ability
var emptyAllyAbilities []*Ability

func scout() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Trade,
				operation: Increase,
				value:     1,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	return &CardEntry{
		cost:                 0,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        emptyAllyAbilities,
	}
}

func viper() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     1,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	return &CardEntry{
		cost:                 0,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        emptyAllyAbilities,
	}
}

func explorer() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Trade,
				operation: Increase,
				value:     2,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	utilizationAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     2,
			}
		},
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 2,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: utilizationAbilities,
		allyAbilities:        emptyAllyAbilities,
	}
}

func blobFighter() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     3,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionRandomCard{
				from: playerDeckMapper(player, Deck),
				to:   playerDeckMapper(player, Hand),
			}
		},
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 1,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        allyAbilities,
	}
}

func tradePod() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Trade,
				operation: Increase,
				value:     3,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     2,
			}
		},
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 2,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        allyAbilities,
	}
}

func ram() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     5,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     2,
			}
		},
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 3,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        allyAbilities,
	}
}

func newMiddleware() *Middleware {
	deck := make(map[string]*CardEntry)
	deck["scout"] = scout()
	deck["viper"] = viper()
	deck["explorer"] = explorer()
	deck["blobFighter"] = blobFighter()
	deck["tradePod"] = tradePod()
	deck["ram"] = ram()

	return &Middleware{
		deck:      deck,
		allyState: emptyAllyState(),
	}
}

func emptyAllyState() *AllyState {
	allyFlags := make(map[Faction]bool)
	allyFlags[Blob] = false
	allyFlags[MachineCult] = false
	allyFlags[StarEmpire] = false
	allyFlags[TradeFederation] = false
	allyAbilitiesState := make(map[Faction]Abilities)
	allyAbilitiesState[Blob] = []*Ability{}
	allyAbilitiesState[MachineCult] = []*Ability{}
	allyAbilitiesState[StarEmpire] = []*Ability{}
	allyAbilitiesState[TradeFederation] = []*Ability{}
	return &AllyState{
		flags:     allyFlags,
		abilities: allyAbilitiesState,
	}
}

func (m *Middleware) resetAllyState() {
	m.allyState = emptyAllyState()
}

func (m *Middleware) handle(action string, player PlayerId) []StateAction {
	currentPlayer := FirstPlayer
	currentPlayerDeck := FirstPlayerDeck
	currentPlayerHand := FirstPlayerHand
	currentPlayerTable := FirstPlayerTable
	currentPlayerDiscard := FirstPlayerDiscard

	opponent := SecondPlayer
	//opponentDeck := SecondPlayerDeck
	//opponentHand := SecondPlayerHand
	//opponentTable := SecondPlayerTable
	//opponentDiscard := SecondPlayerDiscard

	if player == SecondPlayer {
		currentPlayer = SecondPlayer
		currentPlayerDeck = SecondPlayerDeck
		currentPlayerHand = SecondPlayerHand
		currentPlayerTable = SecondPlayerTable
		currentPlayerDiscard = SecondPlayerDiscard

		opponent = FirstPlayer
		//opponentDeck = FirstPlayerDeck
		//opponentHand = FirstPlayerHand
		//opponentTable = FirstPlayerTable
		//opponentDiscard = FirstPlayerDiscard
	}

	var actions []StateAction
	parsed := strings.Split(action, ",")
	if len(parsed) == 0 {
		//TODO: handle exception
		return actions
	}
	action = parsed[0]
	switch action {
	case "play":
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		card, ok := m.deck[strings.Split(id, "_")[0]]
		if !ok {
			//TODO: handle exception
			return actions
		}
		actions = append(actions, &StateActionMoveCard{
			id: id,
			to: currentPlayerTable,
		})
		for _, ability := range card.primaryAbilities {
			if ability.player == Current {
				actions = append(actions, ability.action(currentPlayer))
			} else {
				actions = append(actions, ability.action(opponent))
			}
		}
		allyActivated, ok := m.allyState.flags[card.faction]
		if ok {
			if allyActivated {
				for _, ability := range card.allyAbilities {
					if ability.player == Current {
						actions = append(actions, ability.action(currentPlayer))
					} else {
						actions = append(actions, ability.action(opponent))
					}
				}
				for _, ability := range m.allyState.abilities[card.faction] {
					if ability.player == Current {
						actions = append(actions, ability.action(currentPlayer))
					} else {
						actions = append(actions, ability.action(opponent))
					}
				}
				m.allyState.abilities[card.faction] = []*Ability{}
			} else {
				m.allyState.flags[card.faction] = true
				for _, ability := range card.allyAbilities {
					m.allyState.abilities[card.faction] = append(m.allyState.abilities[card.faction], ability)
				}
			}
		}
	case "end":
		m.resetAllyState()
		actions = append(actions, &StateActionMoveAll{
			from: currentPlayerTable,
			to:   currentPlayerDiscard,
		})
		actions = append(actions, &StateActionChangeCounterValue{
			player:    currentPlayer,
			counter:   Trade,
			operation: Set,
			value:     0,
		})
		actions = append(actions, &StateActionChangeCounterValue{
			player:    currentPlayer,
			counter:   Combat,
			operation: Set,
			value:     0,
		})
		for i := 1; i <= HandCardsQty; i++ {
			actions = append(actions, &StateActionRandomCard{
				from: currentPlayerDeck,
				to:   currentPlayerHand,
			})
		}
		actions = append(actions, &StateActionChangeTurn{})
	case "damage":
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		damage, err := strconv.Atoi(parsed[1])
		if err != nil {
			//TODO: handle exception
			return actions
		}
		actions = append(actions, &StateActionChangeCounterValue{
			player:    opponent,
			counter:   Authority,
			operation: Decrease,
			value:     damage,
		})
		actions = append(actions, &StateActionChangeCounterValue{
			player:    currentPlayer,
			counter:   Combat,
			operation: Decrease,
			value:     damage,
		})
	case "buy":
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		card, ok := m.deck[strings.Split(id, "_")[0]]
		if !ok {
			//TODO: handle exception
			return actions
		}
		actions = append(actions, &StateActionMoveCard{
			id: id,
			to: currentPlayerDiscard,
		})
		actions = append(actions, &StateActionRandomCard{
			from: TradeDeck,
			to:   TradeRow,
		})
		actions = append(actions, &StateActionChangeCounterValue{
			player:    currentPlayer,
			counter:   Trade,
			operation: Decrease,
			value:     card.cost,
		})
	case "utilize":
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		card, ok := m.deck[strings.Split(id, "_")[0]]
		if !ok {
			//TODO: handle exception
			return actions
		}
		actions = append(actions, &StateActionMoveCard{
			id: id,
			to: ScrapHeap,
		})
		for _, ability := range card.utilizationAbilities {
			if ability.player == Current {
				actions = append(actions, ability.action(currentPlayer))
			} else {
				actions = append(actions, ability.action(opponent))
			}
		}
	}
	actions = append(actions, &StateActionGetState{})
	return actions
}

func (m *Middleware) preparareState() []StateAction {
	var actions []StateAction
	for i := 1; i <= FirstPlayerHandCardsQty; i++ {
		actions = append(actions, &StateActionRandomCard{
			from: FirstPlayerDeck,
			to:   FirstPlayerHand,
		})
	}
	for i := 1; i <= HandCardsQty; i++ {
		actions = append(actions, &StateActionRandomCard{
			from: SecondPlayerDeck,
			to:   SecondPlayerHand,
		})
	}
	for i := 1; i <= TradeRowQty; i++ {
		actions = append(actions, &StateActionRandomCard{
			from: TradeDeck,
			to:   TradeRow,
		})
	}
	return actions
}

type DeckType int

const (
	Deck DeckType = iota
	Hand
	Table
	DiscardPile
)

func playerDeckMapper(player PlayerId, deck DeckType) CardLocation {
	type DeckLocation map[DeckType]CardLocation
	type PlayerDeckMapper map[PlayerId]DeckLocation
	firstPlayerDecks := make(map[DeckType]CardLocation)
	firstPlayerDecks[Deck] = FirstPlayerDeck
	firstPlayerDecks[Hand] = FirstPlayerHand
	firstPlayerDecks[Table] = FirstPlayerTable
	firstPlayerDecks[DiscardPile] = FirstPlayerDiscard
	secondPlayerDecks := make(map[DeckType]CardLocation)
	secondPlayerDecks[Deck] = SecondPlayerDeck
	secondPlayerDecks[Hand] = SecondPlayerHand
	secondPlayerDecks[Table] = SecondPlayerTable
	secondPlayerDecks[DiscardPile] = SecondPlayerDiscard
	mapper := make(map[PlayerId]DeckLocation)
	mapper[FirstPlayer] = firstPlayerDecks
	mapper[SecondPlayer] = secondPlayerDecks
	return mapper[player][deck]
}
