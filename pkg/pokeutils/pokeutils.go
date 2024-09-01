package pokeutils

import "math/rand"

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
