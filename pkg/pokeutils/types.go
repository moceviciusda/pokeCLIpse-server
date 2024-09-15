package pokeutils

type IVs struct {
	Hp             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}

type Stats struct {
	Hp             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}

type Move struct {
	Name         string `json:"name"`
	Accuracy     int    `json:"accuracy"`
	Power        int    `json:"power"`
	PP           int    `json:"pp"`
	Type         string `json:"type"`
	DamageClass  string `json:"damage_class"`
	EffectChance int    `json:"effect_chance"`
	Effect       string `json:"effect"`
}

type Pokemon struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Types []string `json:"types"`
	Level int      `json:"level"`
	Shiny bool     `json:"shiny"`
	Stats Stats    `json:"stats"`
	Moves []Move   `json:"moves"`
}
