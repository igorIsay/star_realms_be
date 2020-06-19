package main

type CardEntry struct {
	cost                 int
	qty                  int
	faction              Faction
	primaryAbilities     Abilities
	utilizationAbilities Abilities
	allyAbilities        Abilities
	cardType             CardType
}

func getDeck() map[string]*CardEntry {
	deck := make(map[string]*CardEntry)

	deck["scout"] = scout()
	deck["viper"] = viper()
	deck["explorer"] = explorer()
	deck["blobFighter"] = blobFighter()
	deck["tradePod"] = tradePod()
	deck["ram"] = ram()

	return deck
}

func scout() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Trade,
				operation: Increase,
				value:     1,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	return &CardEntry{
		cost:                 0,
		qty:                  16,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        emptyAllyAbilities,
		cardType:             Ship,
	}
}

func viper() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     1,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	return &CardEntry{
		cost:                 0,
		qty:                  4,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        emptyAllyAbilities,
		cardType:             Ship,
	}
}

func explorer() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Trade,
				operation: Increase,
				value:     2,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	utilizationAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     2,
			}
		},
	}
	utilizationAbilities := []*Ability{&utilizationAbility}
	return &CardEntry{
		cost:                 2,
		qty:                  10,
		faction:              Unaligned,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: utilizationAbilities,
		allyAbilities:        emptyAllyAbilities,
		cardType:             Ship,
	}
}

func blobFighter() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     3,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionRandomCard{
				from: playerDeckMapper(player, Deck),
				to:   playerDeckMapper(player, Hand),
			}
		},
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 1,
		qty:                  3,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
	}
}

func tradePod() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Trade,
				operation: Increase,
				value:     3,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     2,
			}
		},
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 2,
		qty:                  2,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
	}
}

func ram() *CardEntry {
	primaryAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     5,
			}
		},
	}
	primaryAbilities := []*Ability{&primaryAbility}
	allyAbility := Ability{
		player: Current,
		action: func(player PlayerId) StateAction {
			return &StateActionChangeCounterValue{
				player:    player,
				counter:   Combat,
				operation: Increase,
				value:     2,
			}
		},
	}
	allyAbilities := []*Ability{&allyAbility}
	return &CardEntry{
		cost:                 3,
		qty:                  2,
		faction:              Blob,
		primaryAbilities:     primaryAbilities,
		utilizationAbilities: emptyUtilizationAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
	}
}
