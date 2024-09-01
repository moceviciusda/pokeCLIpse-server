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

type Pokemon struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Shiny bool   `json:"shiny"`
	Stats Stats  `json:"stats"`
}
