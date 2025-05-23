package pokeutils

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strings"
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

var TypeIcons = map[string]string{
	"normal":   IconNormal,
	"fighting": IconFighting,
	"flying":   IconFlying,
	"poison":   IconPoison,
	"ground":   IconGround,
	"rock":     IconRock,
	"bug":      IconBug,
	"ghost":    IconGhost,
	"steel":    IconSteel,
	"fire":     IconFire,
	"water":    IconWater,
	"grass":    IconGrass,
	"electric": IconElectric,
	"psychic":  IconPsychic,
	"ice":      IconIce,
	"dragon":   IconDragon,
	"dark":     IconDark,
	"fairy":    IconFairy,
}

var Starters = []string{
	"bulbasaur",
	"charmander",
	"squirtle",
	"chikorita",
	"cyndaquil",
	"totodile",
	"treecko",
	"torchic",
	"mudkip",
	"turtwig",
	"chimchar",
	"piplup",
	"snivy",
	"tepig",
	"oshawott",
	"chespin",
	"fennekin",
	"froakie",
	"rowlet",
	"litten",
	"popplio",
	"grookey",
	"scorbunny",
	"sobble",
	"sprigatito",
	"fuecoco",
	"quaxly",
}

var StarterTypeMap = map[string]string{
	"bulbasaur":  IconGrass + IconPoison,
	"charmander": IconFire,
	"squirtle":   IconWater,
	"chikorita":  IconGrass,
	"cyndaquil":  IconFire,
	"totodile":   IconWater,
	"treecko":    IconGrass,
	"torchic":    IconFire,
	"mudkip":     IconWater,
	"turtwig":    IconGrass,
	"chimchar":   IconFire,
	"piplup":     IconWater,
	"snivy":      IconGrass,
	"tepig":      IconFire,
	"oshawott":   IconWater,
	"chespin":    IconGrass,
	"fennekin":   IconFire,
	"froakie":    IconWater,
	"rowlet":     IconGrass + IconFlying,
	"litten":     IconFire,
	"popplio":    IconWater,
	"grookey":    IconGrass,
	"scorbunny":  IconFire,
	"sobble":     IconWater,
	"sprigatito": IconGrass,
	"fuecoco":    IconFire,
	"quaxly":     IconWater,
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
	levelFactor := (2*float64(attacker.Level))/5 + 2

	var adRatio float64
	if move.DamageClass == "physical" {
		adRatio = float64(attacker.Stats.Attack) / float64(defender.Stats.Defense)
	} else {
		adRatio = float64(attacker.Stats.SpecialAttack) / float64(defender.Stats.SpecialDefense)
	}
	// log.Println("Attack/Defense ratio: ", adRatio)

	baseDamage := (float64(move.Power) * levelFactor * adRatio / 50) + 2
	// log.Println("Damage before modifiers: ", baseDamage)

	var stab float64 = 1
	if slices.Contains(attacker.Types, move.Type) {
		stab = 1.5
	}

	var typeEffectiveness float64 = 1
	for _, t := range defender.Types {
		typeEffectiveness *= TypeEffectiveness(move.Type, t)

	}
	switch typeEffectiveness {
	case 0:
		flavourText = "It doesn't affect " + defender.Name + "..."
	case 0.5:
		flavourText = "It's not very effective..."
	case 2:
		flavourText = "It's super effective!"
	case 4:
		flavourText = "It's ultra effective!"
	}

	damage := baseDamage * stab * typeEffectiveness
	// log.Println("Damage after modifiers: ", damage)

	return int(damage), flavourText
}

func TypeEffectiveness(moveType, targetType string) float64 {
	typeChart := map[string]map[string]float64{
		"normal": {
			"rock":  0.5,
			"ghost": 0,
			"steel": 0.5,
		},
		"fire": {
			"fire":   0.5,
			"water":  0.5,
			"grass":  2,
			"ice":    2,
			"bug":    2,
			"rock":   0.5,
			"dragon": 0.5,
			"steel":  2,
		},
		"water": {
			"fire":   2,
			"water":  0.5,
			"grass":  0.5,
			"ground": 2,
			"rock":   2,
			"dragon": 0.5,
		},
		"electric": {
			"water":    2,
			"electric": 0.5,
			"grass":    0.5,
			"ground":   0,
			"flying":   2,
			"dragon":   0.5,
		},
		"grass": {
			"fire":   0.5,
			"water":  2,
			"grass":  0.5,
			"poison": 0.5,
			"ground": 2,
			"flying": 0.5,
			"bug":    0.5,
			"rock":   2,
			"dragon": 0.5,
			"steel":  0.5,
		},
		"ice": {
			"fire":   0.5,
			"water":  0.5,
			"grass":  2,
			"ice":    0.5,
			"ground": 2,
			"flying": 2,
			"dragon": 2,
			"steel":  0.5,
		},
		"fighting": {
			"normal":  2,
			"ice":     2,
			"poison":  0.5,
			"flying":  0.5,
			"psychic": 0.5,
			"bug":     0.5,
			"rock":    2,
			"ghost":   0,
			"dark":    2,
			"steel":   2,
			"fairy":   0.5,
		},
		"poison": {
			"grass":  2,
			"poison": 0.5,
			"ground": 0.5,
			"rock":   0.5,
			"ghost":  0.5,
			"steel":  0,
			"fairy":  2,
		},
		"ground": {
			"fire":     2,
			"electric": 2,
			"grass":    0.5,
			"poison":   2,
			"flying":   0,
			"bug":      0.5,
			"rock":     2,
			"steel":    2,
		},
		"flying": {
			"electric": 0.5,
			"grass":    2,
			"fighting": 2,
			"bug":      2,
			"rock":     0.5,
			"steel":    0.5,
		},
		"psychic": {
			"fighting": 2,
			"poison":   2,
			"psychic":  0.5,
			"dark":     0,
			"steel":    0.5,
		},
		"bug": {
			"fire":     0.5,
			"grass":    2,
			"fighting": 0.5,
			"poison":   0.5,
			"flying":   0.5,
			"psychic":  2,
			"ghost":    0.5,
			"dark":     2,
			"steel":    0.5,
			"fairy":    0.5,
		},
		"rock": {
			"fire":     2,
			"ice":      2,
			"fighting": 0.5,
			"ground":   0.5,
			"flying":   2,
			"bug":      2,
			"steel":    0.5,
		},
		"ghost": {
			"normal":  0,
			"psychic": 2,
			"ghost":   2,
			"dark":    0.5,
		},
		"dragon": {
			"dragon": 2,
			"steel":  0.5,
			"fairy":  0,
		},
		"dark": {
			"fighting": 0.5,
			"psychic":  2,
			"ghost":    2,
			"dark":     0.5,
			"fairy":    0.5,
		},
		"steel": {
			"fire":     0.5,
			"water":    0.5,
			"electric": 0.5,
			"ice":      2,
			"rock":     2,
			"steel":    0.5,
			"fairy":    2,
		},
		"fairy": {
			"fire":     0.5,
			"fighting": 2,
			"poison":   0.5,
			"dragon":   2,
			"dark":     2,
			"steel":    0.5,
		},
	}

	moveType = strings.ToLower(moveType)
	targetType = strings.ToLower(targetType)

	if effectiveness, ok := typeChart[moveType]; ok {
		if modifier, ok := effectiveness[targetType]; ok {
			return modifier
		}
	}

	return 1.0
}

func ExpAtLevel(level int) int {
	return 4 * (level * level * level) / 5
}

func LevelAtExp(exp int) int {
	return int(math.Cbrt(float64(5*exp) / 4))
}

func ExpYield(baseExp, level int) int {
	return (baseExp * level / 7)
}
