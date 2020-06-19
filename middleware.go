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
			actions = append(actions, &StateActionMoveCard{
				id: id,
				to: currentPlayerTable,
			})
		} else {
			actions = append(actions, &StateActionMoveCard{
				id: id,
				to: currentPlayerBases,
			})
		}
		m.playAbilities(player, card, &actions)
	case End:
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
		actions = append(actions, &StateActionRequestUserAction{
			player: opponent,
			action: Start,
		})
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
		actions = append(actions, &StateActionRequestUserAction{
			player: currentPlayer,
			action: None,
		})
	case DestroyBase:
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
		if card.cardType != Base {
			//TODO: handle exception
			return actions
		}
		actions = append(actions, &StateActionChangeCounterValue{
			player:    currentPlayer,
			counter:   Combat,
			operation: Decrease,
			value:     card.defense,
		})
		actions = append(actions, &StateActionMoveCard{
			id: id,
			to: opponentDiscard,
		})
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
