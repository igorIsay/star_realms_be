package main

type Faction int

const (
	Unaligned Faction = iota
	Blob
	MachineCult
	StarEmpire
	TradeFederation
)

type CardType int

const (
	Ship CardType = iota
	Base
)

type Abilities []*Ability

type Ability struct {
	player PlayerPointer
	action func(PlayerId) StateAction
}

type CardEntry struct {
	cost                 int
	qty                  int
	defense              int
	faction              Faction
	beforePlay           Abilities
	primaryAbilities     Abilities
	utilizationAbilities Abilities
	allyAbilities        Abilities
	cardType             CardType
}

func getDeck() *map[string]*CardEntry {
	deck := make(map[string]*CardEntry)

	deck["scout"] = scout()
	deck["viper"] = viper()
	deck["explorer"] = explorer()
	deck["blobFighter"] = blobFighter()
	deck["tradePod"] = tradePod()
	deck["ram"] = ram()
	deck["theHive"] = theHive()
	deck["blobWheel"] = blobWheel()
	deck["battlePod"] = battlePod()

	deck["corvette"] = corvette()
	deck["dreadnaught"] = dreadnaught()
	deck["imperialFighter"] = imperialFighter()
	deck["imperialFrigate"] = imperialFrigate()
	deck["royalRedoubt"] = royalRedoubt()
	deck["spaceStation"] = spaceStation()
	deck["surveyShip"] = surveyShip()
	deck["warWorld"] = warWorld()

	deck["battleMech"] = battleMech()
	deck["missileBot"] = missileBot()
	deck["supplyBot"] = supplyBot()
	deck["tradeBot"] = tradeBot()

	return &deck
}

func changeCounter(operation Operation, counter Counter, value int) func(PlayerId) StateAction {
	return func(player PlayerId) StateAction {
		return &StateActionChangeCounterValue{
			player:    player,
			counter:   counter,
			operation: operation,
			value:     value,
		}
	}
}

func actionRequest(action UserAction) func(PlayerId) StateAction {
	return func(player PlayerId) StateAction {
		return &StateActionRequestUserAction{
			player: player,
			action: action,
		}
	}
}

func drawCard(player PlayerId) StateAction {
	return &StateActionRandomCard{
		from: playerDeckMapper(player, Deck),
		to:   playerDeckMapper(player, Hand),
	}
}

func scout() *CardEntry {
	return &CardEntry{
		qty:     16,
		faction: Unaligned,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Trade, 1),
			},
		},
		cardType: Ship,
	}
}

func viper() *CardEntry {
	return &CardEntry{
		qty:     4,
		faction: Unaligned,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 1),
			},
		},
		cardType: Ship,
	}
}

func explorer() *CardEntry {
	return &CardEntry{
		cost:    2,
		qty:     10,
		faction: Unaligned,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Trade, 2),
			},
		},
		utilizationAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 2),
			},
		},
		cardType: Ship,
	}
}

func blobFighter() *CardEntry {
	return &CardEntry{
		cost:    1,
		qty:     3,
		faction: Blob,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 3),
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: drawCard,
			},
		},
		cardType: Ship,
	}
}

func tradePod() *CardEntry {
	return &CardEntry{
		cost:    2,
		qty:     2,
		faction: Blob,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Trade, 3),
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 2),
			},
		},
		cardType: Ship,
	}
}

func ram() *CardEntry {
	return &CardEntry{
		cost:    3,
		qty:     2,
		faction: Blob,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 5),
			},
		},
		utilizationAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Trade, 3),
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 2),
			},
		},
		cardType: Ship,
		defense:  0,
	}
}

func theHive() *CardEntry {
	return &CardEntry{
		cost:    5,
		qty:     1,
		faction: Blob,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 3),
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: drawCard,
			},
		},
		cardType: Base,
		defense:  5,
	}
}

func blobWheel() *CardEntry {
	return &CardEntry{
		cost:    3,
		qty:     3,
		faction: Blob,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 1),
			},
		},
		utilizationAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Trade, 3),
			},
		},
		cardType: Base,
		defense:  5,
	}
}

func battlePod() *CardEntry {
	return &CardEntry{
		cost:     2,
		qty:      2,
		faction:  Blob,
		cardType: Ship,
		beforePlay: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCardTradeRow),
			},
		},
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 4),
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 2),
			},
		},
	}
}

func imperialFighter() *CardEntry {
	return &CardEntry{
		cost:    1,
		qty:     3,
		faction: StarEmpire,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 2),
			},
			&Ability{
				player: Opponent,
				action: changeCounter(Increase, Discard, 1),
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 2),
			},
		},
		cardType: Ship,
	}
}

func imperialFrigate() *CardEntry {
	return &CardEntry{
		cost:    1,
		qty:     3,
		faction: StarEmpire,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 4),
			},
			&Ability{
				player: Opponent,
				action: changeCounter(Increase, Discard, 1),
			},
		},
		utilizationAbilities: []*Ability{
			&Ability{
				player: Current,
				action: drawCard,
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 2),
			},
		},
		cardType: Ship,
	}
}

func corvette() *CardEntry {
	return &CardEntry{
		cost:    2,
		qty:     2,
		faction: StarEmpire,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 1),
			},
			&Ability{
				player: Current,
				action: drawCard,
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 2),
			},
		},
		cardType: Ship,
	}
}

func dreadnaught() *CardEntry {
	return &CardEntry{
		cost:    7,
		qty:     1,
		faction: StarEmpire,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 7),
			},
			&Ability{
				player: Current,
				action: drawCard,
			},
		},
		utilizationAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 5),
			},
		},
		cardType: Ship,
	}
}

func royalRedoubt() *CardEntry {
	return &CardEntry{
		cost:     6,
		qty:      1,
		faction:  StarEmpire,
		cardType: Base,
		defense:  6,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 3),
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Opponent,
				action: changeCounter(Increase, Discard, 1),
			},
		},
	}
}

func spaceStation() *CardEntry {
	return &CardEntry{
		cost:     4,
		qty:      2,
		faction:  StarEmpire,
		cardType: Base,
		defense:  4,
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
		utilizationAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Trade, 4),
			},
		},
	}
}

func surveyShip() *CardEntry {
	return &CardEntry{
		cost:     3,
		qty:      3,
		faction:  StarEmpire,
		cardType: Ship,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Trade, 1),
			},
			&Ability{
				player: Current,
				action: drawCard,
			},
		},
		utilizationAbilities: []*Ability{
			&Ability{
				player: Opponent,
				action: changeCounter(Increase, Discard, 1),
			},
		},
	}
}

func warWorld() *CardEntry {
	return &CardEntry{
		cost:     5,
		qty:      1,
		faction:  StarEmpire,
		cardType: Base,
		defense:  4,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 3),
			},
		},
		allyAbilities: []*Ability{
			&Ability{
				player: Current,
				action: changeCounter(Increase, Combat, 4),
			},
		},
	}
}

func battleMech() *CardEntry {
	return &CardEntry{
		cost:     5,
		qty:      1,
		faction:  MachineCult,
		cardType: Ship,
		beforePlay: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCard),
			},
		},
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
	}
}

func missileBot() *CardEntry {
	return &CardEntry{
		cost:     2,
		qty:      3,
		faction:  MachineCult,
		cardType: Ship,
		beforePlay: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCard),
			},
		},
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
}

func supplyBot() *CardEntry {
	return &CardEntry{
		cost:     3,
		qty:      3,
		faction:  MachineCult,
		cardType: Ship,
		beforePlay: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCard),
			},
		},
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
}

func tradeBot() *CardEntry {
	return &CardEntry{
		cost:     1,
		qty:      3,
		faction:  MachineCult,
		cardType: Ship,
		beforePlay: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCard),
			},
		},
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
}
