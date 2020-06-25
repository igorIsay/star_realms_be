package main

import (
	"strconv"
	"strings"
)

type Middleware struct {
	deck      *map[string]*CardEntry
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

type UserAction int

const (
	None UserAction = iota
	Play
	End
	Damage
	Buy
	Utilize
	Start
	DestroyBase
	DiscardCard
	ScrapCard
	ScrapCardAndPlayBattleMech
	ScrapCardAndPlayMissileBot
	ScrapCardAndPlaySupplyBot
	ScrapCardAndPlayTradeBot
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

func (m *Middleware) handle(action string, player PlayerId, state *State) []StateAction {
	var actions []StateAction

	currentPlayer, err := m.relativePlayer(player, Current)
	if err != nil {
		// TODO handle error
		return actions
	}
	opponent, err := m.relativePlayer(player, Opponent)
	if err != nil {
		// TODO handle error
		return actions
	}

	currentPlayerDeck := FirstPlayerDeck
	currentPlayerHand := FirstPlayerHand
	currentPlayerTable := FirstPlayerTable
	currentPlayerDiscard := FirstPlayerDiscard
	currentPlayerBases := FirstPlayerBases

	//opponentDeck := SecondPlayerDeck
	//opponentHand := SecondPlayerHand
	//opponentTable := SecondPlayerTable
	opponentDiscard := SecondPlayerDiscard

	if player == SecondPlayer {
		currentPlayerDeck = SecondPlayerDeck
		currentPlayerHand = SecondPlayerHand
		currentPlayerTable = SecondPlayerTable
		currentPlayerDiscard = SecondPlayerDiscard
		currentPlayerBases = SecondPlayerBases

		opponent = FirstPlayer
		//opponentDeck = FirstPlayerDeck
		//opponentHand = FirstPlayerHand
		//opponentTable = FirstPlayerTable
		opponentDiscard = FirstPlayerDiscard
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
		return actions
	}
	var userAction UserAction
	switch parsedAction {
	case 1:
		userAction = Play
	case 2:
		userAction = End
	case 3:
		userAction = Damage
	case 4:
		userAction = Buy
	case 5:
		userAction = Utilize
	case 6:
		userAction = Start
	case 7:
		userAction = DestroyBase
	case 8:
		userAction = DiscardCard
	case 9:
		userAction = ScrapCard
	case 10:
		userAction = ScrapCardAndPlayBattleMech
	case 11:
		userAction = ScrapCardAndPlayMissileBot
	case 12:
		userAction = ScrapCardAndPlaySupplyBot
	case 13:
		userAction = ScrapCardAndPlayTradeBot
	default:
		//TODO: handle exception
		return actions
	}
	switch userAction {
	case Play:
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
		if card.cardType == Ship {
			m.moveCard(id, currentPlayerTable, &actions)
		} else {
			m.moveCard(id, currentPlayerBases, &actions)
		}
		m.playAbilities(player, card, &actions)
	case End:
		m.resetAllyState()
		m.moveAll(currentPlayerTable, currentPlayerDiscard, &actions)
		m.changeCounterValue(currentPlayer, Set, Trade, 0, &actions)
		m.changeCounterValue(currentPlayer, Set, Combat, 0, &actions)
		for i := 1; i <= HandCardsQty; i++ {
			m.randomCard(currentPlayerDeck, currentPlayerHand, &actions)
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
		m.moveCard(id, currentPlayerDiscard, &actions)
		m.randomCard(TradeDeck, TradeRow, &actions)
		m.changeCounterValue(currentPlayer, Decrease, Trade, card.cost, &actions)
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
			if ability.player == Current {
				actions = append(actions, ability.action(currentPlayer))
			} else {
				actions = append(actions, ability.action(opponent))
			}
		}
	case Start:
		for cardId, card := range state.Cards {
			if card.Location == currentPlayerBases {
				cardEntry, ok := deck[strings.Split(cardId, "_")[0]]
				if !ok {
					//TODO: handle exception
					return actions
				}
				m.playAbilities(player, cardEntry, &actions)
			}
		}
		m.requestUserAction(currentPlayer, None, &actions)
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

		m.moveCard(id, currentPlayerDiscard, &actions)
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
		m.scrapCard(parsed, currentPlayer, &actions)

	case ScrapCardAndPlayBattleMech:
		m.scrapCard(parsed, currentPlayer, &actions)

		battleMech := &CardEntry{
			faction:  MachineCult,
			cardType: Ship,
			primaryAbilities: []*Ability{
				&Ability{
					player: Current,
					action: changeCounter(Increase, Combat, 4),
				},
			},
			allyAbilities: []*Ability{
				&Ability{
					player: Current,
					action: drawCard,
				},
			},
			utilizationAbilities: emptyAbilities,
		}
		m.playAbilities(player, battleMech, &actions)

	case ScrapCardAndPlayMissileBot:
		m.scrapCard(parsed, currentPlayer, &actions)

		missileBot := &CardEntry{
			faction:  MachineCult,
			cardType: Ship,
			primaryAbilities: []*Ability{
				&Ability{
					player: Current,
					action: changeCounter(Increase, Combat, 2),
				},
			},
			allyAbilities: []*Ability{
				&Ability{
					player: Current,
					action: changeCounter(Increase, Combat, 2),
				},
			},
		}
		m.playAbilities(player, missileBot, &actions)

	case ScrapCardAndPlaySupplyBot:
		m.scrapCard(parsed, currentPlayer, &actions)

		supplyBot := &CardEntry{
			faction:  MachineCult,
			cardType: Ship,
			primaryAbilities: []*Ability{
				&Ability{
					player: Current,
					action: changeCounter(Increase, Trade, 2),
				},
			},
			allyAbilities: []*Ability{
				&Ability{
					player: Current,
					action: changeCounter(Increase, Combat, 2),
				},
			},
		}
		m.playAbilities(player, supplyBot, &actions)

	case ScrapCardAndPlayTradeBot:
		m.scrapCard(parsed, currentPlayer, &actions)

		tradeBot := &CardEntry{
			faction:  MachineCult,
			cardType: Ship,
			primaryAbilities: []*Ability{
				&Ability{
					player: Current,
					action: changeCounter(Increase, Trade, 1),
				},
			},
			allyAbilities: []*Ability{
				&Ability{
					player: Current,
					action: changeCounter(Increase, Combat, 2),
				},
			},
		}
		m.playAbilities(player, tradeBot, &actions)
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

type WrongPlayerPointerError struct{}

func (e *WrongPlayerPointerError) Error() string {
	return "wrong PlayerPointer"
}

type WrongPlayerIdError struct{}

func (e *WrongPlayerIdError) Error() string {
	return "wrong PlayerId"
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
			return actualPlayer, &WrongPlayerPointerError{}
		}
	case SecondPlayer:
		switch playerPointer {
		case Current:
			return SecondPlayer, nil
		case Opponent:
			return FirstPlayer, nil
		default:
			return actualPlayer, &WrongPlayerPointerError{}
		}
	default:
		return actualPlayer, &WrongPlayerIdError{}
	}
}

type WrongCountersPointerError struct{}

func (e *WrongCountersPointerError) Error() string {
	return "wrong CountersPointer"
}

func (m *Middleware) relativeCounters(actualPlayer PlayerId, countersPointer CountersPointer, state *State) (Counters, error) {
	switch actualPlayer {
	case FirstPlayer:
		switch countersPointer {
		case CurrentPlayerCounters:
			return state.FirstPlayerCounters, nil
		case OpponentCounters:
			return state.SecondPlayerCounters, nil
		default:
			return Counters{}, &WrongCountersPointerError{}
		}
	case SecondPlayer:
		switch countersPointer {
		case CurrentPlayerCounters:
			return state.SecondPlayerCounters, nil
		case OpponentCounters:
			return state.FirstPlayerCounters, nil
		default:
			return Counters{}, &WrongCountersPointerError{}
		}
	default:
		return Counters{}, &WrongPlayerIdError{}
	}
}

func (m *Middleware) playAbilities(player PlayerId, card *CardEntry, actions *[]StateAction) {
	currentPlayer, err := m.relativePlayer(player, Current)
	if err != nil {
		// TODO handle error
		return
	}
	opponent, err := m.relativePlayer(player, Opponent)
	if err != nil {
		// TODO handle error
		return
	}

	for _, ability := range card.primaryAbilities {
		if ability.player == Current {
			*actions = append(*actions, ability.action(currentPlayer))
		} else {
			*actions = append(*actions, ability.action(opponent))
		}
	}
	allyActivated, ok := m.allyState.flags[card.faction]
	if ok {
		if allyActivated {
			for _, ability := range card.allyAbilities {
				if ability.player == Current {
					*actions = append(*actions, ability.action(currentPlayer))
				} else {
					*actions = append(*actions, ability.action(opponent))
				}
			}
			for _, ability := range m.allyState.abilities[card.faction] {
				if ability.player == Current {
					*actions = append(*actions, ability.action(currentPlayer))
				} else {
					*actions = append(*actions, ability.action(opponent))
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

func (m *Middleware) scrapCard(actionData []string, player PlayerId, actions *[]StateAction) {
	deck := *m.deck
	if len(actionData) > 1 {
		id := actionData[1]
		_, ok := deck[strings.Split(id, "_")[0]]
		if ok {
			m.moveCard(id, ScrapHeap, actions)
		}
	}
	m.requestUserAction(player, None, actions)
}
