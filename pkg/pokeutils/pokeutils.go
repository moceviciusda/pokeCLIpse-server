package pokeutils

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	IconNormal   = "♟️"
	IconFighting = "🥊"
	IconFlying   = "🦅"
	IconPoison   = "☠️"
	IconGround   = "🏜️"
	IconRock     = "⛰️"
	IconBug      = "🐞"
	IconGhost    = "👻"
	IconSteel    = "🔩"
	IconFire     = "🔥"
	IconWater    = "💧"
	IconGrass    = "🌿"
	IconElectric = "⚡"
	IconPsychic  = "🧠"
	IconIce      = "❄️"
	IconDragon   = "🐉"
	IconDark     = "🌙"
	IconFairy    = "✨"
)

var Starters = []string{
	"Bulbasaur",
	"Charmander",
	"Squirtle",
	"Chikorita",
	"Cyndaquil",
	"Totodile",
	"Treecko",
	"Torchic",
	"Mudkip",
	"Turtwig",
	"Chimchar",
	"Piplup",
	"Snivy",
	"Tepig",
	"Oshawott",
	"Chespin",
	"Fennekin",
	"Froakie",
	"Rowlet",
	"Litten",
	"Popplio",
	"Grookey",
	"Scorbunny",
	"Sobble",
	"Sprigatito",
	"Fuecoco",
	"Quaxly",
}

var StarterTypeMap = map[string]string{
	"Bulbasaur":  IconGrass + IconPoison,
	"Charmander": IconFire,
	"Squirtle":   IconWater,
	"Chikorita":  IconGrass,
	"Cyndaquil":  IconFire,
	"Totodile":   IconWater,
	"Treecko":    IconGrass,
	"Torchic":    IconFire,
	"Mudkip":     IconWater,
	"Turtwig":    IconGrass,
	"Chimchar":   IconFire,
	"Piplup":     IconWater,
	"Snivy":      IconGrass,
	"Tepig":      IconFire,
	"Oshawott":   IconWater,
	"Chespin":    IconGrass,
	"Fennekin":   IconFire,
	"Froakie":    IconWater,
	"Rowlet":     IconGrass + IconFlying,
	"Litten":     IconFire,
	"Popplio":    IconWater,
	"Grookey":    IconGrass,
	"Scorbunny":  IconFire,
	"Sobble":     IconWater,
	"Sprigatito": IconGrass,
	"Fuecoco":    IconFire,
	"Quaxly":     IconWater,
}

func GenerateIVs() IVs {
	return IVs{
		Hp:             rand.Intn(32),
		Attack:         rand.Intn(32),
		Defense:        rand.Intn(32),
		SpecialAttack:  rand.Intn(32),
		SpecialDefense: rand.Intn(32),
		Speed:          rand.Intn(32),
	}
}

func CalculateStat(baseStat, iv, level int) int {
	return ((2*baseStat + iv) * level / 100) + 5
}

func IsShiny() bool {
	return rand.Intn(128) == 0
}

func CalculateStats(baseStats Stats, ivs IVs, level int) Stats {
	return Stats{
		Hp:             CalculateStat(baseStats.Hp, ivs.Hp, level) + 5,
		Attack:         CalculateStat(baseStats.Attack, ivs.Attack, level),
		Defense:        CalculateStat(baseStats.Defense, ivs.Defense, level),
		SpecialAttack:  CalculateStat(baseStats.SpecialAttack, ivs.SpecialAttack, level),
		SpecialDefense: CalculateStat(baseStats.SpecialDefense, ivs.SpecialDefense, level),
		Speed:          CalculateStat(baseStats.Speed, ivs.Speed, level),
	}
}

func (s Stats) String() string {
	return fmt.Sprintf(
		`	HP: %d			Speed: %d
	Attack: %d		Special Attack: %d
	Defense: %d		Special Defense: %d`, s.Hp, s.Speed, s.Attack, s.SpecialAttack, s.Defense, s.SpecialDefense)
}

func CalculateDamage(attacker, defender Pokemon, move Move) (dmg int, flavourText string) {
	var ad float64
	if move.DamageClass == "physical" {
		ad = float64(attacker.Stats.Attack) / float64(defender.Stats.Defense)
	} else {
		ad = float64(attacker.Stats.SpecialAttack) / float64(defender.Stats.SpecialDefense)
	}
	log.Println("Attack/Defense ratio: ", ad)

	damage := ((float64(attacker.Level)*2/5+2)*float64(move.Power)*ad/50 + 2)
	log.Println("Damage before modifiers: ", damage)

	var stab float64 = 1
	for _, t := range attacker.Types {
		if t == move.Type {
			stab = 1.5
			break
		}
	}

	var typeEffectiveness float64 = 1
	for _, t := range defender.Types {
		te := TypeEffectiveness(move.Type, t)
		typeEffectiveness *= te
		switch te {
		case 0:
			flavourText = "It doesn't affect " + defender.Name + "..."
		case 0.5:
			flavourText = "It's not very effective..."
		case 2:
			flavourText = "It's super effective!"
		}
	}

	damage = damage * stab * typeEffectiveness
	log.Println("Damage after modifiers: ", damage)

	return int(damage), flavourText
}

func TypeEffectiveness(moveType, targetType string) float64 {
	typeChart := map[string]map[string]float64{
		"Normal": {
			"Rock":  0.5,
			"Ghost": 0,
			"Steel": 0.5,
		},
		"Fire": {
			"Fire":   0.5,
			"Water":  0.5,
			"Grass":  2,
			"Ice":    2,
			"Bug":    2,
			"Rock":   0.5,
			"Dragon": 0.5,
			"Steel":  2,
		},
		"Water": {
			"Fire":   2,
			"Water":  0.5,
			"Grass":  0.5,
			"Ground": 2,
			"Rock":   2,
			"Dragon": 0.5,
		},
		"Electric": {
			"Water":    2,
			"Electric": 0.5,
			"Grass":    0.5,
			"Ground":   0,
			"Flying":   2,
			"Dragon":   0.5,
		},
		"Grass": {
			"Fire":   0.5,
			"Water":  2,
			"Grass":  0.5,
			"Poison": 0.5,
			"Ground": 2,
			"Flying": 0.5,
			"Bug":    0.5,
			"Rock":   2,
			"Dragon": 0.5,
			"Steel":  0.5,
		},
		"Ice": {
			"Fire":   0.5,
			"Water":  0.5,
			"Grass":  2,
			"Ice":    0.5,
			"Ground": 2,
			"Flying": 2,
			"Dragon": 2,
			"Steel":  0.5,
		},
		"Fighting": {
			"Normal":  2,
			"Ice":     2,
			"Poison":  0.5,
			"Flying":  0.5,
			"Psychic": 0.5,
			"Bug":     0.5,
			"Rock":    2,
			"Ghost":   0,
			"Dark":    2,
			"Steel":   2,
			"Fairy":   0.5,
		},
		"Poison": {
			"Grass":  2,
			"Poison": 0.5,
			"Ground": 0.5,
			"Rock":   0.5,
			"Ghost":  0.5,
			"Steel":  0,
			"Fairy":  2,
		},
		"Ground": {
			"Fire":     2,
			"Electric": 2,
			"Grass":    0.5,
			"Poison":   2,
			"Flying":   0,
			"Bug":      0.5,
			"Rock":     2,
			"Steel":    2,
		},
		"Flying": {
			"Electric": 0.5,
			"Grass":    2,
			"Fighting": 2,
			"Bug":      2,
			"Rock":     0.5,
			"Steel":    0.5,
		},
		"Psychic": {
			"Fighting": 2,
			"Poison":   2,
			"Psychic":  0.5,
			"Dark":     0,
			"Steel":    0.5,
		},
		"Bug": {
			"Fire":     0.5,
			"Grass":    2,
			"Fighting": 0.5,
			"Poison":   0.5,
			"Flying":   0.5,
			"Psychic":  2,
			"Ghost":    0.5,
			"Dark":     2,
			"Steel":    0.5,
			"Fairy":    0.5,
		},
		"Rock": {
			"Fire":     2,
			"Ice":      2,
			"Fighting": 0.5,
			"Ground":   0.5,
			"Flying":   2,
			"Bug":      2,
			"Steel":    0.5,
		},
		"Ghost": {
			"Normal":  0,
			"Psychic": 2,
			"Ghost":   2,
			"Dark":    0.5,
		},
		"Dragon": {
			"Dragon": 2,
			"Steel":  0.5,
			"Fairy":  0,
		},
		"Dark": {
			"Fighting": 0.5,
			"Psychic":  2,
			"Ghost":    2,
			"Dark":     0.5,
			"Fairy":    0.5,
		},
		"Steel": {
			"Fire":     0.5,
			"Water":    0.5,
			"Electric": 0.5,
			"Ice":      2,
			"Rock":     2,
			"Steel":    0.5,
			"Fairy":    2,
		},
		"Fairy": {
			"Fire":     0.5,
			"Fighting": 2,
			"Poison":   0.5,
			"Dragon":   2,
			"Dark":     2,
			"Steel":    0.5,
		},
	}

	if effectiveness, ok := typeChart[moveType]; ok {
		if modifier, ok := effectiveness[targetType]; ok {
			return modifier
		}
	}

	return 1.0
}
