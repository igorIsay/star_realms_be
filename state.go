package main

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
	SecondPlayerDeck
	SecondPlayerHand
	SecondPlayerTable
	SecondPlayerDiscard
)

type Counters struct {
	Trade     int `json:"trade"`
	Combat    int `json:"combat"`
	Authority int `json:"authority"`
	Discard   int `json:"discard"`
}

func newState() *State {
	const initialAuthority int = 50
	cards := cardsInitialSet()
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

func cardsInitialSet() map[string]*Card {
	cards := make(map[string]*Card)
	firstPlayerDeck := [...]string{"scout_1", "scout_2", "scout_3", "scout_4", "scout_5", "scout_6", "scout_7", "scout_8", "viper_1", "viper_2"}
	secondPlayerDeck := [...]string{"scout_9", "scout_10", "scout_11", "scout_12", "scout_13", "scout_14", "scout_15", "scout_16", "viper_3", "viper_4"}
	explorersDeck := [...]string{"explorer_1", "explorer_2", "explorer_3", "explorer_4", "explorer_5", "explorer_6", "explorer_7", "explorer_8", "explorer_9", "explorer_10"}
	blobFighters := [...]string{"blobFighter_1", "blobFighter_2", "blobFighter_3"}
	tradePods := [...]string{"tradePod_1", "tradePod_2"}
	rams := [...]string{"ram_1", "ram_2"}
	for _, id := range firstPlayerDeck {
		cards[id] = &Card{
			Location: FirstPlayerDeck,
		}
	}
	for _, id := range secondPlayerDeck {
		cards[id] = &Card{
			Location: SecondPlayerDeck,
		}
	}
	for _, id := range explorersDeck {
		cards[id] = &Card{
			Location: Explorers,
		}
	}
	for _, id := range blobFighters {
		cards[id] = &Card{
			Location: TradeDeck,
		}
	}
	for _, id := range tradePods {
		cards[id] = &Card{
			Location: TradeDeck,
		}
	}
	for _, id := range rams {
		cards[id] = &Card{
			Location: TradeDeck,
		}
	}
	return cards
}
