package main

import "fmt"

type State struct {
	Turn                      PlayerId                      `json:"turn"`
	FirstPlayerCounters       Counters                      `json:"firstPlayerCounters"`
	SecondPlayerCounters      Counters                      `json:"secondPlayerCounters"`
	Cards                     map[string]*Card              `json:"cards"`
	FirstPlayerActionRequest  ActionRequest                 `json:"firstPlayerActionRequest"`
	SecondPlayerActionRequest ActionRequest                 `json:"secondPlayerActionRequest"`
	ActivatedAbilities        map[string]ActivatedAbilities `json:"activatedAbilities"`
	lastIndex                 map[CardLocation]int
}

type ActivatedAbilities map[AbilityId]bool

type Card struct {
	Location CardLocation `json:"location"`
	index    int
}

type CardLocation int

const (
	UndefinedLocation CardLocation = iota
	TradeDeck
	TradeRow
	Explorers
	ScrapHeap
	FirstPlayerDeck
	FirstPlayerHand
	FirstPlayerTable
	FirstPlayerDiscard
	FirstPlayerBases
	SecondPlayerDeck
	SecondPlayerHand
	SecondPlayerTable
	SecondPlayerDiscard
	SecondPlayerBases
)

type Counters struct {
	Trade     int `json:"trade"`
	Combat    int `json:"combat"`
	Authority int `json:"authority"`
	Discard   int `json:"discard"`
}

type UserAction int

const (
	NoneAction UserAction = iota
	Play
	End
	Damage
	Buy
	Utilize
	Start
	DestroyBase
	DiscardCard
	ActivateAbility
	ScrapCard
	ScrapCardTradeRow
	DestroyBaseMissileMech
)

type ActionRequest struct {
	Action UserAction `json:"action"`
	CardId string     `json:"cardId"`
}

func newState(deck *map[string]*CardEntry) *State {
	const initialAuthority int = 50
	lastIndex := make(map[CardLocation]int)
	cards := cardsInitialSet(deck, lastIndex)
	return &State{
		Turn: FirstPlayer,
		FirstPlayerCounters: Counters{
			Authority: initialAuthority,
		},
		SecondPlayerCounters: Counters{
			Authority: initialAuthority,
		},
		Cards: cards,
		FirstPlayerActionRequest:  ActionRequest{},
		SecondPlayerActionRequest: ActionRequest{},
		ActivatedAbilities:        make(map[string]ActivatedAbilities),
		lastIndex:                 lastIndex,
	}
}

func cardsInitialSet(deck *map[string]*CardEntry, lastIndex map[CardLocation]int) map[string]*Card {
	cards := make(map[string]*Card)
	for key, card := range *deck {
		switch key {
		case "scout":
			h := card.qty / 2
			for i := 1; i <= card.qty; i++ {
				id := fmt.Sprintf("%s_%d", key, i)
				if i <= h {
					lastIndex[FirstPlayerDeck] += 1
					cards[id] = &Card{
						Location: FirstPlayerDeck,
						index:    lastIndex[FirstPlayerDeck],
					}
				} else {
					lastIndex[SecondPlayerDeck] += 1
					cards[id] = &Card{
						Location: SecondPlayerDeck,
						index:    lastIndex[SecondPlayerDeck],
					}
				}
			}
		case "viper":
			h := card.qty / 2
			for i := 1; i <= card.qty; i++ {
				id := fmt.Sprintf("%s_%d", key, i)
				if i <= h {
					lastIndex[FirstPlayerDeck] += 1
					cards[id] = &Card{
						Location: FirstPlayerDeck,
						index:    lastIndex[FirstPlayerDeck],
					}
				} else {
					lastIndex[SecondPlayerDeck] += 1
					cards[id] = &Card{
						Location: SecondPlayerDeck,
						index:    lastIndex[SecondPlayerDeck],
					}
				}
			}
		case "explorer":
			for i := 1; i <= card.qty; i++ {
				lastIndex[Explorers] += 1
				id := fmt.Sprintf("%s_%d", key, i)
				cards[id] = &Card{
					Location: Explorers,
					index:    lastIndex[Explorers],
				}
			}
		default:
			lastIndex[TradeDeck] += 1
			for i := 1; i <= card.qty; i++ {
				id := fmt.Sprintf("%s_%d", key, i)
				cards[id] = &Card{
					Location: TradeDeck,
					index:    lastIndex[TradeDeck],
				}
			}

		}
	}
	return cards
}
