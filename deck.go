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
	primaryAbilities     Abilities
	utilizationAbilities Abilities
	allyAbilities        Abilities
	cardType             CardType
}

var emptyAbilities []*Ability

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
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Trade, 1),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	return &CardEntry{
		cost:                 0,
		qty:                  16,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyAbilities,
		allyAbilities:        emptyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func viper() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 1),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	return &CardEntry{
		cost:                 0,
		qty:                  4,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyAbilities,
		allyAbilities:        emptyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func explorer() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Trade, 2),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	utilizationAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 2,
		qty:                  10,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: utilizationAbilities,
		allyAbilities:        emptyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func blobFighter() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 3),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: drawCard,
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 1,
		qty:                  3,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func tradePod() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Trade, 3),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 2,
		qty:                  2,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func ram() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 5),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	allyAbilities := []*Ability{&allyAbility}
	utilizationAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Trade, 3),
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 3,
		qty:                  2,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: utilizationAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func theHive() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 3),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: drawCard,
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 5,
		qty:                  1,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Base,
		defense:              5,
	}
}

func blobWheel() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 1),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	utilizationAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Trade, 3),
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 3,
		qty:                  3,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: utilizationAbilities,
		allyAbilities:        emptyAbilities,
		cardType:             Base,
		defense:              5,
	}
}

func imperialFighter() *CardEntry {
	primaryAbilityCombat := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	primaryAbilityDiscard := Ability{
		player: Opponent,
		action: changeCounter(Increase, Discard, 1),
	}
	primaryAbilities := []*Ability{
		&primaryAbilityCombat,
		&primaryAbilityDiscard,
	}
	allyAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 1,
		qty:                  3,
		faction:              StarEmpire,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func imperialFrigate() *CardEntry {
	primaryAbilityCombat := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 4),
	}
	primaryAbilityDiscard := Ability{
		player: Opponent,
		action: changeCounter(Increase, Discard, 1),
	}
	primaryAbilities := []*Ability{
		&primaryAbilityCombat,
		&primaryAbilityDiscard,
	}
	allyAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	allyAbilities := []*Ability{&allyAbility}
	utilizationAbility := Ability{
		player: Current,
		action: drawCard,
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 1,
		qty:                  3,
		faction:              StarEmpire,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: utilizationAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func corvette() *CardEntry {
	primaryAbilityCombat := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 1),
	}
	primaryAbilityDraw := Ability{
		player: Current,
		action: drawCard,
	}
	primaryAbilities := []*Ability{
		&primaryAbilityCombat,
		&primaryAbilityDraw,
	}
	allyAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 2,
		qty:                  2,
		faction:              StarEmpire,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func dreadnaught() *CardEntry {
	primaryAbilityCombat := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 7),
	}
	primaryAbilityDraw := Ability{
		player: Current,
		action: drawCard,
	}
	primaryAbilities := []*Ability{
		&primaryAbilityCombat,
		&primaryAbilityDraw,
	}
	utilizationAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 5),
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 7,
		qty:                  1,
		faction:              StarEmpire,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: utilizationAbilities,
		allyAbilities:        emptyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func royalRedoubt() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 3),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Opponent,
		action: changeCounter(Increase, Discard, 1),
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 6,
		qty:                  1,
		faction:              StarEmpire,
		cardType:             Base,
		defense:              6,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
	}
}

func spaceStation() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 2),
	}
	allyAbilities := []*Ability{&allyAbility}
	utilizationAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Trade, 4),
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 4,
		qty:                  2,
		faction:              StarEmpire,
		cardType:             Base,
		defense:              4,
		primaryAbilities:     primaryAbilities,
		allyAbilities:        allyAbilities,
		utilizationAbilities: utilizationAbilities,
	}
}

func surveyShip() *CardEntry {
	primaryAbilityTrade := Ability{
		player: Current,
		action: changeCounter(Increase, Trade, 1),
	}
	primaryAbilityDraw := Ability{
		player: Current,
		action: drawCard,
	}
	primaryAbilities := []*Ability{
		&primaryAbilityTrade,
		&primaryAbilityDraw,
	}
	utilizationAbility := Ability{
		player: Opponent,
		action: changeCounter(Increase, Discard, 1),
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 3,
		qty:                  3,
		faction:              StarEmpire,
		cardType:             Ship,
		defense:              0,
		primaryAbilities:     primaryAbilities,
		allyAbilities:        emptyAbilities,
		utilizationAbilities: utilizationAbilities,
	}
}

func warWorld() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 3),
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: changeCounter(Increase, Combat, 4),
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 5,
		qty:                  1,
		faction:              StarEmpire,
		cardType:             Base,
		defense:              4,
		primaryAbilities:     primaryAbilities,
		allyAbilities:        allyAbilities,
		utilizationAbilities: emptyAbilities,
	}
}

func battleMech() *CardEntry {
	return &CardEntry{
		cost:     5,
		qty:      1,
		faction:  Unaligned,
		cardType: Ship,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCardAndPlayBattleMech),
			},
		},
	}
}

func missileBot() *CardEntry {
	return &CardEntry{
		cost:     2,
		qty:      3,
		faction:  Unaligned,
		cardType: Ship,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCardAndPlayMissileBot),
			},
		},
	}
}

func supplyBot() *CardEntry {
	return &CardEntry{
		cost:     3,
		qty:      3,
		faction:  Unaligned,
		cardType: Ship,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCardAndPlaySupplyBot),
			},
		},
	}
}

func tradeBot() *CardEntry {
	return &CardEntry{
		cost:     1,
		qty:      3,
		faction:  Unaligned,
		cardType: Ship,
		primaryAbilities: []*Ability{
			&Ability{
				player: Current,
				action: actionRequest(ScrapCardAndPlayTradeBot),
			},
		},
	}
}
