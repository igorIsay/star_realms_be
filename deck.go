package main

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

	return &deck
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
		utilizationAbilities: emptyAbilities,
		allyAbilities:        emptyAbilities,
		cardType:             Ship,
		defense:              0,
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
		utilizationAbilities: emptyAbilities,
		allyAbilities:        emptyAbilities,
		cardType:             Ship,
		defense:              0,
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
		allyAbilities:        emptyAbilities,
		cardType:             Ship,
		defense:              0,
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
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
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
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
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
		utilizationAbilities: emptyAbilities,
		allyAbilities:        allyAbilities,
		cardType:             Ship,
		defense:              0,
	}
}

func theHive() *CardEntry {
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
	utilizationAbility := Ability{
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
