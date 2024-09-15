package pokebattle

import (
	"log"
	"math/rand"
	"time"

	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
)

type Trainer struct {
	Name    string
	Pokemon []pokeutils.Pokemon
}

type PokemonState struct {
	HP, Attack, Defense, SpecialAttack, SpecialDefense, Speed int
}

type Battle struct {
	Trainers [2]Trainer
}

func NewBattle(trainer1, trainer2 Trainer) *Battle {
	return &Battle{
		Trainers: [2]Trainer{trainer1, trainer2},
	}
}

func (b *Battle) Run() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for i := 0; i < 2; i++ {
			log.Println()
			attacker := &b.Trainers[i]
			defender := &b.Trainers[(i+1)%2]

			move, ok := selectMove(&attacker.Pokemon[0])
			if !ok {
				log.Printf("%s's %s has no moves left!\n", attacker.Name, attacker.Pokemon[0].Name)
				continue
			}

			if move.Accuracy < rand.Intn(100) {
				log.Printf("%s's %s used %s but missed!\n", attacker.Name, attacker.Pokemon[0].Name, move.Name)
				continue
			}

			damage, flavourText := pokeutils.CalculateDamage(attacker.Pokemon[0], defender.Pokemon[0], move)
			log.Printf("%s's %s used %s and dealt %d damage to %s's %s\n", attacker.Name, attacker.Pokemon[0].Name, move.Name, damage, defender.Name, defender.Pokemon[0].Name)
			if flavourText != "" {
				log.Println(flavourText)
			}

			defender.Pokemon[0].Stats.Hp -= damage

			if defender.Pokemon[0].Stats.Hp <= 0 {
				log.Println(defender.Name + "'s " + defender.Pokemon[0].Name + " fainted!")
				return
			}
		}

	}
}

func selectMove(p *pokeutils.Pokemon) (pokeutils.Move, bool) {
	validMoves := make([]pokeutils.Move, 0, len(p.Moves))
	for _, move := range p.Moves {
		if move.PP > 0 {
			validMoves = append(validMoves, move)
		}
	}
	if len(validMoves) == 0 {
		log.Println(p.Name + " has no moves left!")
		return pokeutils.Move{}, false
	}

	move := validMoves[rand.Intn(len(validMoves))]
	move.PP--
	return move, true
}
