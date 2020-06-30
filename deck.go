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

type AbilityActionType int

const (
	Instant AbilityActionType = iota
	Activated
)

type AbilityId int

const (
	DefaultAbility AbilityId = iota
	Utilization
	PatrolMechTrade
	PatrolMechCombat
	PatrolMechScrap
	BlobCarrierAcquire
	BlobDestroyerDestroyBase
)

type AbilityGroup int

const (
	Primary AbilityGroup = iota
	Ally
	BeforePlay
)

type Abilities []*Ability

type Ability struct {
	group      AbilityGroup
	actionType AbilityActionType
	id         AbilityId
	player     PlayerPointer
	actions    func(PlayerId, string) []StateAction
}

type CardEntry struct {
	cost       int
	qty        int
	defense    int
	faction    Faction
	beforePlay Abilities
	abilities  Abilities
	cardType   CardType
}

func getDeck() *map[string]*CardEntry {
	deck := make(map[string]*CardEntry)

	deck["scout"] = scout()
	deck["viper"] = viper()
	deck["explorer"] = explorer()

	deck["blobFighter"] = blobFighter()
	deck["tradePod"] = tradePod()
	deck["ram"] = ram()
	deck["battlePod"] = battlePod()
	deck["theHive"] = theHive()
	deck["blobWheel"] = blobWheel()
	deck["blobCarrier"] = blobCarrier()
	deck["blobDestroyer"] = blobDestroyer()

	deck["corvette"] = corvette()
	deck["dreadnaught"] = dreadnaught()
	deck["imperialFighter"] = imperialFighter()
	deck["imperialFrigate"] = imperialFrigate()
	deck["royalRedoubt"] = royalRedoubt()
	deck["spaceStation"] = spaceStation()
	deck["surveyShip"] = surveyShip()
	deck["warWorld"] = warWorld()
	deck["battlecruiser"] = battlecruiser()

	deck["battleMech"] = battleMech()
	deck["missileBot"] = missileBot()
	deck["supplyBot"] = supplyBot()
	deck["missileMech"] = missileMech()
	deck["tradeBot"] = tradeBot()
	deck["patrolMech"] = patrolMech()

	return &deck
}

func changeCounter(operation Operation, counter Counter, value int) func(PlayerId, string) []StateAction {
	return func(player PlayerId, cardId string) []StateAction {
		return []StateAction{
			&StateActionChangeCounterValue{
				player:    player,
				counter:   counter,
				operation: operation,
				value:     value,
			},
		}
	}
}

func actionRequest(action UserAction) func(PlayerId, string) []StateAction {
	return func(player PlayerId, cardId string) []StateAction {
		return []StateAction{
			&StateActionRequestUserAction{
				player: player,
				action: action,
				cardId: cardId,
			},
		}
	}
}

func drawCard(player PlayerId, cardId string) []StateAction {
	return []StateAction{
		&StateActionTopCard{
			from: playerDeckMapper(player, Deck),
			to:   playerDeckMapper(player, Hand),
		},
	}
}

func scout() *CardEntry {
	return &CardEntry{
		qty:     16,
		faction: Unaligned,
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Trade, 1),
			},
		},
		cardType: Ship,
	}
}

func viper() *CardEntry {
	return &CardEntry{
		qty:     4,
		faction: Unaligned,
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 1),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Trade, 2),
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         Utilization,
				player:     Current,
				actions:    changeCounter(Increase, Combat, 2),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 3),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: drawCard,
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Trade, 3),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 5),
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         Utilization,
				player:     Current,
				actions:    changeCounter(Increase, Trade, 3),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 3),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: drawCard,
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 1),
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         Utilization,
				player:     Current,
				actions:    changeCounter(Increase, Trade, 3),
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
				group:   BeforePlay,
				player:  Current,
				actions: actionRequest(ScrapCardTradeRow),
			},
		},
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 4),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
			},
		},
	}
}

func blobCarrier() *CardEntry {
	return &CardEntry{
		cost:     6,
		qty:      1,
		faction:  Blob,
		cardType: Ship,
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 7),
			},
			&Ability{
				group:      Ally,
				actionType: Activated,
				id:         BlobCarrierAcquire,
				player:     Current,
				actions:    actionRequest(AcquireShipForFree),
			},
		},
	}
}

func blobDestroyer() *CardEntry {
	return &CardEntry{
		cost:     4,
		qty:      2,
		faction:  Blob,
		cardType: Ship,
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 6),
			},
			&Ability{
				group:      Ally,
				actionType: Activated,
				id:         BlobDestroyerDestroyBase,
				player:     Current,
				actions:    actionRequest(DestroyBaseBlobDestroyer),
			},
		},
	}
}

func imperialFighter() *CardEntry {
	return &CardEntry{
		cost:    1,
		qty:     3,
		faction: StarEmpire,
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
			},
			&Ability{
				group:   Primary,
				player:  Opponent,
				actions: changeCounter(Increase, Discard, 1),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 4),
			},
			&Ability{
				group:   Primary,
				player:  Opponent,
				actions: changeCounter(Increase, Discard, 1),
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         Utilization,
				player:     Current,
				actions:    drawCard,
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 1),
			},
			&Ability{
				group:   Primary,
				player:  Current,
				actions: drawCard,
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 7),
			},
			&Ability{
				group:   Primary,
				player:  Current,
				actions: drawCard,
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         Utilization,
				player:     Current,
				actions:    changeCounter(Increase, Combat, 5),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 3),
			},
			&Ability{
				group:   Ally,
				player:  Opponent,
				actions: changeCounter(Increase, Discard, 1),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         Utilization,
				player:     Current,
				actions:    changeCounter(Increase, Trade, 4),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Trade, 1),
			},
			&Ability{
				group:   Primary,
				player:  Current,
				actions: drawCard,
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         Utilization,
				player:     Opponent,
				actions:    changeCounter(Increase, Discard, 1),
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
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 3),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 4),
			},
		},
	}
}

func battlecruiser() *CardEntry {
	return &CardEntry{
		cost:     6,
		qty:      1,
		faction:  StarEmpire,
		cardType: Ship,
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 5),
			},
			&Ability{
				group:   Primary,
				player:  Current,
				actions: drawCard,
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         Utilization,
				player:     Current,
				actions: func(player PlayerId, cardId string) []StateAction {
					return []StateAction{
						&StateActionTopCard{
							from: playerDeckMapper(player, Deck),
							to:   playerDeckMapper(player, Hand),
						},
						&StateActionRequestUserAction{
							player: player,
							action: DestroyBaseForFree,
							cardId: cardId,
						},
					}
				},
			},
			&Ability{
				group:   Ally,
				player:  Opponent,
				actions: changeCounter(Increase, Discard, 1),
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
				group:   BeforePlay,
				player:  Current,
				actions: actionRequest(ScrapCard),
			},
		},
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 4),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: drawCard,
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
				group:   BeforePlay,
				player:  Current,
				actions: actionRequest(ScrapCard),
			},
		},
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
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
				group:   BeforePlay,
				player:  Current,
				actions: actionRequest(ScrapCard),
			},
		},
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Trade, 2),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
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
				group:   BeforePlay,
				player:  Current,
				actions: actionRequest(ScrapCard),
			},
		},
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Trade, 1),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: changeCounter(Increase, Combat, 2),
			},
		},
	}
}

func missileMech() *CardEntry {
	return &CardEntry{
		cost:     6,
		qty:      1,
		faction:  MachineCult,
		cardType: Ship,
		beforePlay: []*Ability{
			&Ability{
				group:   BeforePlay,
				player:  Current,
				actions: actionRequest(DestroyBaseForFree),
			},
		},
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: changeCounter(Increase, Combat, 6),
			},
			&Ability{
				group:   Ally,
				player:  Current,
				actions: drawCard,
			},
		},
	}
}

func patrolMech() *CardEntry {
	return &CardEntry{
		cost:     4,
		qty:      2,
		faction:  MachineCult,
		cardType: Ship,
		abilities: []*Ability{
			&Ability{
				group:   Primary,
				player:  Current,
				actions: actionRequest(ActivateAbility),
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         PatrolMechTrade,
				player:     Current,
				actions: func(player PlayerId, cardId string) []StateAction {
					return []StateAction{
						&StateActionChangeCounterValue{
							player:    player,
							counter:   Trade,
							operation: Increase,
							value:     3,
						},
						&StateActionDisableActivatedAbility{
							cardId:    cardId,
							abilityId: PatrolMechCombat,
						},
						&StateActionRequestUserAction{
							player: player,
							action: NoneAction,
						},
					}
				},
			},
			&Ability{
				group:      Primary,
				actionType: Activated,
				id:         PatrolMechCombat,
				player:     Current,
				actions: func(player PlayerId, cardId string) []StateAction {
					return []StateAction{
						&StateActionChangeCounterValue{
							player:    player,
							counter:   Combat,
							operation: Increase,
							value:     5,
						},
						&StateActionDisableActivatedAbility{
							cardId:    cardId,
							abilityId: PatrolMechTrade,
						},
						&StateActionRequestUserAction{
							player: player,
							action: NoneAction,
						},
					}
				},
			},
			&Ability{
				group:      Ally,
				actionType: Activated,
				id:         PatrolMechScrap,
				player:     Current,
				actions:    actionRequest(ScrapCard),
			},
		},
	}
}
