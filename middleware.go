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

type ActionRequestPointer int

const (
	CurrentPlayerActionRequest ActionRequestPointer = iota
	OpponentActionRequest
)

const NEEDLE_SUFFIX string = "_needle"
const NEEDLE_ID string = "stealthNeedle_1"

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

func (m *Middleware) activateAbility(ability *Ability, cardId string, player PlayerId, state *State, actions *[]StateAction) {
	currentPlayer, err := playerByPointer(player, Current)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return
	}
	opponent, err := playerByPointer(player, Opponent)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return
	}

	var abilityActions []StateAction
	if ability.player == Current {
		abilityActions = ability.actions(currentPlayer, cardId, state)
	} else {
		abilityActions = ability.actions(opponent, cardId, state)
	}
	for _, action := range abilityActions {
		*actions = append(*actions, action)
	}
}

func (m *Middleware) processAbility(ability *Ability, cardId string, player PlayerId, state *State, actions *[]StateAction) {
	switch ability.actionType {
	case Instant:
		m.activateAbility(ability, cardId, player, state, actions)
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
	actions = append(actions, &StateActionResetActions{})

	currentPlayer, err := playerByPointer(player, Current)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return actions
	}
	currentPlayerCounters, err := countersByPointer(player, CurrentPlayerCounters, state)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return actions
	}
	currentPlayerActionRequest, err := actionRequestByPointer(player, CurrentPlayerActionRequest, state)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return actions
	}
	opponent, err := playerByPointer(player, Opponent)
	if err != nil {
		// TODO handle error
		log.Println(err)
		return actions
	}

	currentDeck, err := locationByPointer(CurrentDeck, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	currentHand, err := locationByPointer(CurrentHand, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	currentTable, err := locationByPointer(CurrentTable, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	currentBases, err := locationByPointer(CurrentBases, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	currentDiscard, err := locationByPointer(CurrentDiscard, player)
	if err != nil {
		// TODO: handle exception
		log.Println(err)
		return actions
	}
	opponentDiscard, err := locationByPointer(OpponentDiscard, player)
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

	userAction := UserAction(parsedAction)

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
				m.moveCard(id, state.Cards[id].Location, currentTable, &actions)
				if currentPlayerCounters.fleetFlag == 1 {
					m.changeCounterValue(currentPlayer, Increase, Combat, 1, &actions)
				}
			} else {
				m.moveCard(id, state.Cards[id].Location, currentBases, &actions)
			}

			if card.faction == Blob {
				m.changeCounterValue(currentPlayer, Increase, blobs, 1, &actions)
			}
			if len(card.beforePlay) > 0 {
				for _, ability := range card.beforePlay {
					m.processAbility(ability, id, player, state, &actions)
				}

				m.deferredCall = func() []StateAction {
					var actions []StateAction
					m.playAbilities(player, id, state, &actions)
					return actions
				}
			} else {
				m.playAbilities(player, id, state, &actions)
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

			abilityId := AbilityId(parsedAbilityId)
			for _, ability := range card.abilities {
				if ability.id == abilityId {
					m.activateAbility(ability, id, player, state, &actions)

					if ability.id == Utilization {
						if strings.HasSuffix(id, NEEDLE_SUFFIX) {
							m.moveCard(NEEDLE_ID, state.Cards[NEEDLE_ID].Location, ScrapHeap, &actions)
						} else {
							m.moveCard(id, state.Cards[id].Location, ScrapHeap, &actions)
						}
						// Update AllyState after utilization
						foundSameFactionCard := false
						for cardId, c := range state.Cards {
							cardEntryId := strings.Split(cardId, "_")[0]
							if cardId != id &&
								(c.Location == currentTable || c.Location == currentBases) &&
								(deck[cardEntryId].faction == card.faction || cardEntryId == "mechWorld") {

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
		m.changeCounterValue(currentPlayer, Set, ShipsOnTop, 0, &actions)
		m.changeCounterValue(currentPlayer, Set, fleetFlag, 0, &actions)
		m.changeCounterValue(currentPlayer, Set, blobs, 0, &actions)
		for i := 1; i <= HandCardsQty; i++ {
			m.topCard(currentDeck, currentHand, &actions)
		}
		opponentCounters, err := countersByPointer(player, OpponentCounters, state)
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
		cardEntryId := strings.Split(id, "_")[0]
		card, ok := deck[cardEntryId]
		if !ok {
			//TODO: handle exception
			return actions
		}

		if card.cardType == Ship && currentPlayerCounters.ShipsOnTop > 0 {
			m.moveCard(id, state.Cards[id].Location, currentDeck, &actions)
			m.changeCounterValue(currentPlayer, Decrease, ShipsOnTop, 1, &actions)
		} else {
			m.moveCard(id, state.Cards[id].Location, currentDiscard, &actions)
		}
		m.changeCounterValue(currentPlayer, Decrease, Trade, card.cost, &actions)
		if cardEntryId != "explorer" {
			m.topCard(TradeDeck, TradeRow, &actions)
		}
	case Start:
		for cardId, card := range state.Cards {
			if card.Location == currentBases {
				m.playAbilities(player, cardId, state, &actions)
			}
		}
		actionRequested := false
		for _, action := range actions {
			if action.Type() == RequestUserAction {
				actionRequested = true
			}
		}
		if !actionRequested {
			m.requestUserAction(currentPlayer, NoneAction, &actions)
		}
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
		if card.cardType == Ship {
			//TODO: handle exception
			return actions
		}
		m.changeCounterValue(currentPlayer, Decrease, Combat, card.defense, &actions)
		m.moveCard(baseId, state.Cards[baseId].Location, opponentDiscard, &actions)
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

		m.moveCard(id, state.Cards[id].Location, currentDiscard, &actions)
		m.changeCounterValue(currentPlayer, Decrease, Discard, 1, &actions)

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
				m.moveCard(id, state.Cards[id].Location, ScrapHeap, &actions)
			}
		}
		m.requestUserAction(player, NoneAction, &actions)
	case ScrapCardTradeRow:
		if len(parsed) > 1 {
			id := parsed[1]
			_, ok := deck[strings.Split(id, "_")[0]]
			if ok {
				m.moveCard(id, state.Cards[id].Location, ScrapHeap, &actions)
			}
			m.topCard(TradeDeck, TradeRow, &actions)
		}
		m.requestUserAction(player, NoneAction, &actions)
	case ScrapCardInHand:
		if len(parsed) > 1 {
			id := parsed[1]
			_, ok := deck[strings.Split(id, "_")[0]]
			if ok {
				card, ok := state.Cards[id]
				if ok && card.Location == currentHand {
					m.moveCard(id, state.Cards[id].Location, ScrapHeap, &actions)
					m.requestUserAction(player, NoneAction, &actions)
				}
			}
		}
	case DestroyBaseForFree:
		if len(parsed) > 1 {
			baseId := parsed[1]
			card, ok := deck[strings.Split(baseId, "_")[0]]
			if !ok {
				//TODO: handle exception
				return actions
			}
			if card.cardType == Ship {
				//TODO: handle exception
				return actions
			}
			m.moveCard(baseId, state.Cards[baseId].Location, opponentDiscard, &actions)
		}
		m.requestUserAction(player, NoneAction, &actions)
	case DestroyBaseBlobDestroyer:
		if len(parsed) > 1 {
			baseId := parsed[1]
			card, ok := deck[strings.Split(baseId, "_")[0]]
			if !ok {
				//TODO: handle exception
				return actions
			}
			if card.cardType == Ship {
				//TODO: handle exception
				return actions
			}
			m.moveCard(baseId, state.Cards[baseId].Location, opponentDiscard, &actions)
		}
		m.requestUserAction(player, ScrapCardTradeRow, &actions)
	case AcquireShipForFree:
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		cardEntryId := strings.Split(id, "_")[0]
		card, ok := deck[cardEntryId]
		if !ok {
			//TODO: handle exception
			return actions
		}
		if card.cardType != Ship {
			//TODO: handle exception
			return actions
		}
		m.moveCard(id, state.Cards[id].Location, currentDeck, &actions)
		if cardEntryId != "explorer" {
			m.topCard(TradeDeck, TradeRow, &actions)
		}
		m.requestUserAction(player, NoneAction, &actions)
	case ActivateBrainWorld:
		if currentPlayerActionRequest.Action == ActivateBrainWorld {
			if len(parsed) > 1 && len(parsed) <= 3 {
				for _, id := range parsed[1:] {
					_, ok := deck[strings.Split(id, "_")[0]]
					if ok {
						card, ok := state.Cards[id]
						if ok && (card.Location == currentHand || card.Location == currentDiscard) {
							m.moveCard(id, state.Cards[id].Location, ScrapHeap, &actions)
							m.topCard(currentDeck, currentHand, &actions)
						}
					}
				}
			}
			m.requestUserAction(player, NoneAction, &actions)
		}
	case ActivateRecyclingStation:
		if currentPlayerActionRequest.Action == ActivateRecyclingStation {
			if len(parsed) > 1 && len(parsed) <= 3 {
				for _, id := range parsed[1:] {
					_, ok := deck[strings.Split(id, "_")[0]]
					if ok {
						card, ok := state.Cards[id]
						if ok && (card.Location == currentHand) {
							m.moveCard(id, state.Cards[id].Location, currentDiscard, &actions)
							m.topCard(currentDeck, currentHand, &actions)
						}
					}
				}
			}
			m.requestUserAction(player, NoneAction, &actions)
		}
	case ActivateMechWorld:
		if currentPlayerActionRequest.Action == ActivateMechWorld {
			for faction := range m.allyState.flags {
				m.allyState.flags[faction] = true
			}
			for faction, abilities := range m.allyState.abilities {
				for _, cardAbility := range abilities {
					m.processAbility(cardAbility.ability, cardAbility.cardId, player, state, &actions)
				}
				m.allyState.abilities[faction] = []*CardAbility{}
			}
			m.requestUserAction(player, NoneAction, &actions)
		}
	case ActivateNeedle:
		if len(parsed) < 2 {
			//TODO: handle exception
			return actions
		}
		id := parsed[1]
		card, ok := state.Cards[id]
		if !ok {
			//TODO: handle exception
			return actions
		}
		if card.Location != currentTable {
			//TODO: handle exception
			return actions
		}
		cardEntryId := strings.Split(id, "_")[0]
		cardEntry, ok := deck[cardEntryId]
		if !ok {
			//TODO: handle exception
			return actions
		}
		if cardEntry.cardType != Ship {
			//TODO: handle exception
			return actions
		}
		if len(cardEntry.beforePlay) > 0 {
			for _, ability := range cardEntry.beforePlay {
				m.processAbility(ability, id, player, state, &actions)
			}

			m.deferredCall = func() []StateAction {
				var actions []StateAction
				m.playAbilities(player, id+NEEDLE_SUFFIX, state, &actions)
				return actions
			}
		} else {
			m.playAbilities(player, id+NEEDLE_SUFFIX, state, &actions)
		}
		actionRequested := false
		for _, action := range actions {
			if action.Type() == RequestUserAction {
				actionRequested = true
			}
		}
		if !actionRequested {
			m.requestUserAction(currentPlayer, NoneAction, &actions)
		}
	}

	for _, action := range deferredActions {
		actions = append(actions, action)
	}

	actions = append(actions, &StateActionGetState{})
	return actions
}

func (m *Middleware) prepareState() []StateAction {
	var actions []StateAction
	m.shuffle(TradeDeck, &actions)
	m.shuffle(FirstPlayerDeck, &actions)
	m.shuffle(SecondPlayerDeck, &actions)
	for i := 1; i <= FirstPlayerHandCardsQty; i++ {
		m.topCard(FirstPlayerDeck, FirstPlayerHand, &actions)
	}
	for i := 1; i <= HandCardsQty; i++ {
		m.topCard(SecondPlayerDeck, SecondPlayerHand, &actions)
	}
	for i := 1; i <= TradeRowQty; i++ {
		m.topCard(TradeDeck, TradeRow, &actions)
	}
	actions = append(actions, &StateActionResetActions{})
	return actions
}

func (m *Middleware) playAbilities(player PlayerId, cardId string, state *State, actions *[]StateAction) {
	deck := *m.deck
	card, ok := deck[strings.Split(cardId, "_")[0]]
	if ok {
		allyAbilities := []*CardAbility{}
		for _, ability := range card.abilities {
			switch ability.group {
			case Primary:
				m.processAbility(ability, cardId, player, state, actions)
			case Ally:
				allyAbilities = append(allyAbilities, &CardAbility{ability: ability, cardId: cardId})
			}
		}
		allyActivated, ok := m.allyState.flags[card.faction]
		if ok {
			if allyActivated {
				for _, cardAbility := range allyAbilities {
					m.processAbility(cardAbility.ability, cardAbility.cardId, player, state, actions)
				}
				for _, cardAbility := range m.allyState.abilities[card.faction] {
					m.processAbility(cardAbility.ability, cardAbility.cardId, player, state, actions)
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

func (m *Middleware) moveCard(id string, from CardLocation, to CardLocation, actions *[]StateAction) {
	*actions = append(*actions, &StateActionMoveCard{
		id:   id,
		to:   to,
		from: from,
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

func (m *Middleware) topCard(from CardLocation, to CardLocation, actions *[]StateAction) {
	*actions = append(*actions, &StateActionTopCard{
		from: from,
		to:   to,
	})
}

func (m *Middleware) shuffle(deck CardLocation, actions *[]StateAction) {
	*actions = append(*actions, &StateActionShuffleDeck{
		deck: deck,
	})
}

func (m *Middleware) moveAll(from CardLocation, to CardLocation, actions *[]StateAction) {
	*actions = append(*actions, &StateActionMoveAll{
		from: from,
		to:   to,
	})
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

type WrongActionRequestPointerError struct {
	p ActionRequestPointer
}

func (e *WrongActionRequestPointerError) Error() string {
	return fmt.Sprintf("wrong ActionRequestPointer %d", e.p)
}

type WrongLocationPointerError struct {
	p LocationPointer
}

func (e *WrongLocationPointerError) Error() string {
	return fmt.Sprintf("wrong LocationPointer %d", e.p)
}

func playerByPointer(actualPlayer PlayerId, playerPointer PlayerPointer) (PlayerId, error) {
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

func locationByPointer(pointer LocationPointer, player PlayerId) (CardLocation, error) {
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

func countersByPointer(player PlayerId, countersPointer CountersPointer, state *State) (Counters, error) {
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

func actionRequestByPointer(player PlayerId, pointer ActionRequestPointer, state *State) (ActionRequest, error) {
	switch player {
	case FirstPlayer:
		switch pointer {
		case CurrentPlayerActionRequest:
			return state.FirstPlayerActionRequest, nil
		case OpponentActionRequest:
			return state.SecondPlayerActionRequest, nil
		default:
			return ActionRequest{}, &WrongActionRequestPointerError{pointer}
		}
	case SecondPlayer:
		switch pointer {
		case CurrentPlayerActionRequest:
			return state.SecondPlayerActionRequest, nil
		case OpponentActionRequest:
			return state.FirstPlayerActionRequest, nil
		default:
			return ActionRequest{}, &WrongActionRequestPointerError{pointer}
		}
	default:
		return ActionRequest{}, &WrongPlayerIdError{player}
	}
}
