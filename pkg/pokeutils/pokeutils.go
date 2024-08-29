package pokeutils

import "math/rand"

type IVs struct {
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}

func GenerateIVs() IVs {
	return IVs{
		HP:             rand.Intn(32),
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
