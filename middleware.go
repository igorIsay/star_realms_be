package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Middleware struct {
	deck         *map[string]*CardEntry
	allyState    *AllyState
	deferredCall func() []StateAction
}

type AllyState struct {
	flags     map[Faction]bool
	abilities map[Faction]CardAbilities
}

type CardAbilities []*CardAbility

type CardAbility struct {
	cardId  string
	ability *Ability
}

const (
	TradeRowQty             int = 5
	HandCardsQty            int = 5
	FirstPlayerHandCardsQty int = 3
)

type PlayerPointer int

const (
	Current PlayerPointer = iota
	Opponent
)

type LocationPointer int

const (
	CurrentDeck LocationPointer = iota
	CurrentHand
	CurrentTable
	CurrentDiscard
	CurrentBases
	OpponentDeck
	OpponentHand
	OpponentTable
	OpponentDiscard
	OpponentBases
)

type CountersPointer int

const (
	CurrentPlayerCounters CountersPointer = iota
	OpponentCounters
)

func newMiddleware(deck *map[string]*CardEntry) *Middleware {
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
	allyAbilitiesState := make(map[Faction]CardAbilities)
	allyAbilitiesState[Blob] = []*CardAbility{}
	allyAbilitiesState[MachineCult] = []*CardAbility{}
	allyAbilitiesState[StarEmpire] = []*CardAbility{}
	allyAbilitiesState[TradeFederation] = []*CardAbility{}
	return &AllyState{
		flags:     allyFlags,
		abilities: allyAbilitiesState,
	}
}

func (m *Middleware) resetAllyState() {
	m.allyState = emptyAllyState()
}

func (m *Middleware) activateAbility(ability *Ability, cardId string, player PlayerId, actions *[]StateAction) {
	currentPlayer, err := m.relativePlayer(player, Current)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return
	}
	opponent, err := m.relativePlayer(player, Opponent)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return
	}

	var abilityActions []StateAction
	if ability.player == Current {
		abilityActions = ability.actions(currentPlayer, cardId)
	} else {
		abilityActions = ability.actions(opponent, cardId)
	}
	for _, action := range abilityActions {
		*actions = append(*actions, action)
	}
}

func (m *Middleware) processAbility(ability *Ability, cardId string, player PlayerId, actions *[]StateAction) {
	switch ability.actionType {
	case Instant:
		m.activateAbility(ability, cardId, player, actions)
	case Activated:
		*actions = append(
			*actions,
			&StateActionAddActivatedAbility{
				cardId:    cardId,
				abilityId: ability.id,
			},
		)
	}
}

func (m *Middleware) handle(action string, player PlayerId, state *State) []StateAction {
	var actions []StateAction
	var deferredActions []StateAction

	currentPlayer, err := m.relativePlayer(player, Current)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return actions
	}
	opponent, err := m.relativePlayer(player, Opponent)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return actions
	}

	currentDeck, err := m.locationByPointer(CurrentDeck, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	currentHand, err := m.locationByPointer(CurrentHand, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	currentTable, err := m.locationByPointer(CurrentTable, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	currentBases, err := m.locationByPointer(CurrentBases, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	currentDiscard, err := m.locationByPointer(CurrentDiscard, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	opponentDiscard, err := m.locationByPointer(OpponentDiscard, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}

	deck := *m.deck
	parsed := strings.Split(action, ",")
	if len(parsed) == 0 {
		//TODO: handle exception
		return actions
	}
	parsedAction, err := strconv.Atoi(parsed[0])
	if err != nil {
		//TODO: handle exception
		log.Println(err)
		return actions
	}
	var userAction UserAction
	switch parsedAction {
	case int(Play):
		userAction = Play
	case int(End):
		userAction = End
	case int(Damage):
		userAction = Damage
	case int(Buy):
		userAction = Buy
	case int(Utilize):
		userAction = Utilize
	case int(Start):
		userAction = Start
	case int(DestroyBase):
		userAction = DestroyBase
	case int(DiscardCard):
		userAction = DiscardCard
	case int(ActivateAbility):
		userAction = ActivateAbility
	case int(ScrapCard):
		userAction = ScrapCard
	case int(ScrapCardTradeRow):
		userAction = ScrapCardTradeRow
	case int(DestroyBaseMissileMech):
		userAction = DestroyBaseMissileMech
	default:
		//TODO: handle exception
		return actions
	}
	if m.deferredCall != nil {
		deferredActions = m.deferredCall()
		m.deferredCall = nil
	}
	switch userAction {
	case Play:
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		card, ok := deck[strings.Split(id, "_")[0]]
		if ok {
			if card.cardType == Ship {
				m.moveCard(id, currentTable, &actions)
			} else {
				m.moveCard(id, currentBases, &actions)
			}

			if len(card.beforePlay) > 0 {
				for _, ability := range card.beforePlay {
					m.processAbility(ability, id, player, &actions)
				}

				m.deferredCall = func() []StateAction {
					var actions []StateAction
					m.playAbilities(player, id, &actions)
					return actions
				}
			} else {
				m.playAbilities(player, id, &actions)
			}
		}

	case ActivateAbility:
		if len(parsed) < 3 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		card, ok := deck[strings.Split(id, "_")[0]]
		if ok {
			parsedAbilityId, err := strconv.Atoi(parsed[2])
			if err != nil {
				//TODO: handle exception
				return actions
			}

			var abilityId AbilityId
			switch parsedAbilityId {
			case int(Utilization):
				abilityId = Utilization
			case int(PatrolMechTrade):
				abilityId = PatrolMechTrade
			case int(PatrolMechCombat):
				abilityId = PatrolMechCombat
			case int(PatrolMechScrap):
				abilityId = PatrolMechScrap
			default:
				return actions
			}

			for _, ability := range card.abilities {
				if ability.id == abilityId {
					m.activateAbility(ability, id, player, &actions)

					if ability.id == Utilization {
						m.moveCard(id, ScrapHeap, &actions)
						// Update AllyState after utilization
						foundSameFactionCard := false
						for cardId, c := range state.Cards {
							if cardId != id &&
								(c.Location == currentTable || c.Location == currentBases) &&
								deck[strings.Split(cardId, "_")[0]].faction == card.faction {

								foundSameFactionCard = true
							}
						}
						if foundSameFactionCard == false {
							m.allyState.flags[card.faction] = false
							m.allyState.abilities[card.faction] = []*CardAbility{}
						}
					}
				}
			}
			actions = append(
				actions,
				&StateActionDisableActivatedAbility{
					cardId:    id,
					abilityId: abilityId,
				},
			)
		}

	case End:
		m.resetAllyState()
		m.moveAll(currentTable, currentDiscard, &actions)
		m.changeCounterValue(currentPlayer, Set, Trade, 0, &actions)
		m.changeCounterValue(currentPlayer, Set, Combat, 0, &actions)
		for i := 1; i <= HandCardsQty; i++ {
			m.randomCard(currentDeck, currentHand, &actions)
		}
		opponentCounters, err := m.relativeCounters(player, OpponentCounters, state)
		if err != nil {
			// TODO handle exception
			return actions
		}
		if opponentCounters.Discard > 0 {
			m.requestUserAction(opponent, DiscardCard, &actions)
		} else {
			m.requestUserAction(opponent, Start, &actions)
		}
		actions = append(actions, &StateActionResetActivatedAbilities{})
		actions = append(actions, &StateActionChangeTurn{})
	case Damage:
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		damage, err := strconv.Atoi(parsed[1])
		if err != nil {
			//TODO: handle exception
			return actions
		}
		m.changeCounterValue(opponent, Decrease, Authority, damage, &actions)
		m.changeCounterValue(currentPlayer, Decrease, Combat, damage, &actions)
	case Buy:
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		card, ok := deck[strings.Split(id, "_")[0]]
		if !ok {
			//TODO: handle exception
			return actions
		}
		m.moveCard(id, currentDiscard, &actions)
		m.randomCard(TradeDeck, TradeRow, &actions)
		m.changeCounterValue(currentPlayer, Decrease, Trade, card.cost, &actions)
		/*
			case Utilize:
				if len(parsed) < 2 {
					//TODO: handle exception
					return actions
				}
				id := parsed[1]
				card, ok := deck[strings.Split(id, "_")[0]]
				if !ok {
					//TODO: handle exception
					return actions
				}
				m.moveCard(id, ScrapHeap, &actions)
				for _, ability := range card.utilizationAbilities {
					m.processAbility(ability, id, player, &actions)
				}

				// Upade AllyState after utilization
				foundSameFactionCard := false
				for cardId, c := range state.Cards {
					if cardId != id &&
						(c.Location == currentTable || c.Location == currentBases) &&
						deck[strings.Split(cardId, "_")[0]].faction == card.faction {

						foundSameFactionCard = true
					}
				}
				if foundSameFactionCard == false {
					m.allyState.flags[card.faction] = false
					m.allyState.abilities[card.faction] = []*Ability{}
				}
		*/
	case Start:
		for cardId, card := range state.Cards {
			if card.Location == currentBases {
				m.playAbilities(player, cardId, &actions)
			}
		}
		m.requestUserAction(currentPlayer, NoneAction, &actions)
	case DestroyBase:
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		baseId := parsed[1]
		card, ok := deck[strings.Split(baseId, "_")[0]]
		if !ok {
			//TODO: handle exception
			return actions
		}
		if card.cardType != Base {
			//TODO: handle exception
			return actions
		}
		m.changeCounterValue(currentPlayer, Decrease, Combat, card.defense, &actions)
		m.moveCard(baseId, opponentDiscard, &actions)
	case DiscardCard:
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		_, ok := deck[strings.Split(id, "_")[0]]
		if !ok {
			//TODO: handle exception
			return actions
		}

		m.moveCard(id, currentDiscard, &actions)
		m.changeCounterValue(currentPlayer, Decrease, Discard, 1, &actions)

		currentPlayerCounters, err := m.relativeCounters(player, CurrentPlayerCounters, state)
		if err != nil {
			// TODO handle exception
			return actions
		}
		if currentPlayerCounters.Discard == 1 {
			m.requestUserAction(currentPlayer, Start, &actions)
		}
	case ScrapCard:
		if len(parsed) > 1 {
			id := parsed[1]
			_, ok := deck[strings.Split(id, "_")[0]]
			if ok {
				m.moveCard(id, ScrapHeap, &actions)
			}
		}
		m.requestUserAction(player, NoneAction, &actions)
	case ScrapCardTradeRow:
		if len(parsed) > 1 {
			id := parsed[1]
			_, ok := deck[strings.Split(id, "_")[0]]
			if ok {
				m.moveCard(id, ScrapHeap, &actions)
			}
		}
		m.randomCard(TradeDeck, TradeRow, &actions)
		m.requestUserAction(player, NoneAction, &actions)
	case DestroyBaseMissileMech:
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		baseId := parsed[1]
		card, ok := deck[strings.Split(baseId, "_")[0]]
		if !ok {
			//TODO: handle exception
			return actions
		}
		if card.cardType != Base {
			//TODO: handle exception
			return actions
		}
		m.moveCard(baseId, opponentDiscard, &actions)
	}
	for _, action := range deferredActions {
		actions = append(actions, action)
	}

	actions = append(actions, &StateActionGetState{})
	return actions
}

func (m *Middleware) preparareState() []StateAction {
	var actions []StateAction
	for i := 1; i <= FirstPlayerHandCardsQty; i++ {
		m.randomCard(FirstPlayerDeck, FirstPlayerHand, &actions)
	}
	for i := 1; i <= HandCardsQty; i++ {
		m.randomCard(SecondPlayerDeck, SecondPlayerHand, &actions)
	}
	for i := 1; i <= TradeRowQty; i++ {
		m.randomCard(TradeDeck, TradeRow, &actions)
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

type WrongPlayerIdError struct {
	id PlayerId
}

func (e *WrongPlayerIdError) Error() string {
	return fmt.Sprintf("wrong PlayerId %d", e.id)
}

type WrongPlayerPointerError struct {
	p PlayerPointer
}

func (e *WrongPlayerPointerError) Error() string {
	return fmt.Sprintf("wrong PlayerPointer %d", e.p)
}

type WrongCountersPointerError struct {
	p CountersPointer
}

func (e *WrongCountersPointerError) Error() string {
	return fmt.Sprintf("wrong CountersPointer %d", e.p)
}

type WrongLocationPointerError struct {
	p LocationPointer
}

func (e *WrongLocationPointerError) Error() string {
	return fmt.Sprintf("wrong LocationPointer %d", e.p)
}

func (m *Middleware) relativePlayer(actualPlayer PlayerId, playerPointer PlayerPointer) (PlayerId, error) {
	switch actualPlayer {
	case FirstPlayer:
		switch playerPointer {
		case Current:
			return FirstPlayer, nil
		case Opponent:
			return SecondPlayer, nil
		default:
			return actualPlayer, &WrongPlayerPointerError{playerPointer}
		}
	case SecondPlayer:
		switch playerPointer {
		case Current:
			return SecondPlayer, nil
		case Opponent:
			return FirstPlayer, nil
		default:
			return actualPlayer, &WrongPlayerPointerError{playerPointer}
		}
	default:
		return actualPlayer, &WrongPlayerIdError{actualPlayer}
	}
}

func (m *Middleware) relativeCounters(player PlayerId, countersPointer CountersPointer, state *State) (Counters, error) {
	switch player {
	case FirstPlayer:
		switch countersPointer {
		case CurrentPlayerCounters:
			return state.FirstPlayerCounters, nil
		case OpponentCounters:
			return state.SecondPlayerCounters, nil
		default:
			return Counters{}, &WrongCountersPointerError{countersPointer}
		}
	case SecondPlayer:
		switch countersPointer {
		case CurrentPlayerCounters:
			return state.SecondPlayerCounters, nil
		case OpponentCounters:
			return state.FirstPlayerCounters, nil
		default:
			return Counters{}, &WrongCountersPointerError{countersPointer}
		}
	default:
		return Counters{}, &WrongPlayerIdError{player}
	}
}

func (m *Middleware) playAbilities(player PlayerId, cardId string, actions *[]StateAction) {
	deck := *m.deck
	card, ok := deck[strings.Split(cardId, "_")[0]]
	if ok {
		allyAbilities := []*CardAbility{}
		for _, ability := range card.abilities {
			switch ability.group {
			case Primary:
				m.processAbility(ability, cardId, player, actions)
			case Ally:
				allyAbilities = append(allyAbilities, &CardAbility{ability: ability, cardId: cardId})
			}
		}
		allyActivated, ok := m.allyState.flags[card.faction]
		if ok {
			if allyActivated {
				for _, cardAbility := range allyAbilities {
					m.processAbility(cardAbility.ability, cardAbility.cardId, player, actions)
				}
				for _, cardAbility := range m.allyState.abilities[card.faction] {
					m.processAbility(cardAbility.ability, cardAbility.cardId, player, actions)
				}
				m.allyState.abilities[card.faction] = []*CardAbility{}
			} else {
				m.allyState.flags[card.faction] = true
				for _, cardAbility := range allyAbilities {
					m.allyState.abilities[card.faction] = append(m.allyState.abilities[card.faction], cardAbility)
				}
			}
		}
	}
}

func (m *Middleware) moveCard(id string, to CardLocation, actions *[]StateAction) {
	*actions = append(*actions, &StateActionMoveCard{
		id: id,
		to: to,
	})
}

func (m *Middleware) changeCounterValue(player PlayerId, operation Operation, counter Counter, value int, actions *[]StateAction) {
	*actions = append(*actions, &StateActionChangeCounterValue{
		player:    player,
		counter:   counter,
		operation: operation,
		value:     value,
	})
}

func (m *Middleware) requestUserAction(player PlayerId, action UserAction, actions *[]StateAction) {
	*actions = append(*actions, &StateActionRequestUserAction{
		player: player,
		action: action,
	})
}

func (m *Middleware) randomCard(from CardLocation, to CardLocation, actions *[]StateAction) {
	*actions = append(*actions, &StateActionRandomCard{
		from: from,
		to:   to,
	})
}

func (m *Middleware) moveAll(from CardLocation, to CardLocation, actions *[]StateAction) {
	*actions = append(*actions, &StateActionMoveAll{
		from: from,
		to:   to,
	})
}

func (m *Middleware) locationByPointer(pointer LocationPointer, player PlayerId) (CardLocation, error) {
	switch player {
	case FirstPlayer:
		switch pointer {
		case CurrentTable:
			return FirstPlayerTable, nil
		case CurrentBases:
			return FirstPlayerBases, nil
		case CurrentDeck:
			return FirstPlayerDeck, nil
		case CurrentHand:
			return FirstPlayerHand, nil
		case CurrentDiscard:
			return FirstPlayerDiscard, nil
		case OpponentDiscard:
			return SecondPlayerDiscard, nil
		default:
			return UndefinedLocation, &WrongLocationPointerError{pointer}
		}
	case SecondPlayer:
		switch pointer {
		case CurrentTable:
			return SecondPlayerTable, nil
		case CurrentBases:
			return SecondPlayerBases, nil
		case CurrentDeck:
			return SecondPlayerDeck, nil
		case CurrentHand:
			return SecondPlayerHand, nil
		case CurrentDiscard:
			return SecondPlayerDiscard, nil
		case OpponentDiscard:
			return FirstPlayerDiscard, nil
		default:
			return UndefinedLocation, &WrongLocationPointerError{pointer}
		}
	default:
		return UndefinedLocation, &WrongPlayerIdError{player}
	}
}
