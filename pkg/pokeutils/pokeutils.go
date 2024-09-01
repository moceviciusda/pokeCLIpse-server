package pokeutils

import (
	"fmt"
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
