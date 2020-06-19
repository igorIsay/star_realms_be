package main

import "fmt"

type State struct {
	Turn                 PlayerId         `json:"turn"`
	FirstPlayerCounters  Counters         `json:"firstPlayerCounters"`
	SecondPlayerCounters Counters         `json:"secondPlayerCounters"`
	Cards                map[string]*Card `json:"cards"`
}

type Card struct {
	Location CardLocation
}

type CardLocation int

const (
	TradeDeck CardLocation = iota
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

func newState(deck *map[string]*CardEntry) *State {
	const initialAuthority int = 50
	cards := cardsInitialSet(deck)
	return &State{
		Turn: FirstPlayer,
		FirstPlayerCounters: Counters{
			Authority: initialAuthority,
		},
		SecondPlayerCounters: Counters{
			Authority: initialAuthority,
		},
		Cards: cards,
	}
}

func cardsInitialSet(deck *map[string]*CardEntry) map[string]*Card {
	cards := make(map[string]*Card)
	for key, card := range *deck {
		switch key {
		case "scout":
			h := card.qty / 2
			for i := 1; i <= card.qty; i++ {
				id := fmt.Sprintf("%s_%d", key, i)
				if i <= h {
					cards[id] = &Card{
						Location: FirstPlayerDeck,
					}
				} else {
					cards[id] = &Card{
						Location: SecondPlayerDeck,
					}
				}
			}
		case "viper":
			h := card.qty / 2
			for i := 1; i <= card.qty; i++ {
				id := fmt.Sprintf("%s_%d", key, i)
				if i <= h {
					cards[id] = &Card{
						Location: FirstPlayerDeck,
					}
				} else {
					cards[id] = &Card{
						Location: SecondPlayerDeck,
					}
				}
			}
		case "explorer":
			for i := 1; i <= card.qty; i++ {
				id := fmt.Sprintf("%s_%d", key, i)
				cards[id] = &Card{
					Location: Explorers,
				}
			}
		default:
			for i := 1; i <= card.qty; i++ {
				id := fmt.Sprintf("%s_%d", key, i)
				cards[id] = &Card{
					Location: TradeDeck,
				}
			}

		}
	}
	return cards
}
