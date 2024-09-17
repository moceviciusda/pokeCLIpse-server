package pokebattle

import (
	"fmt"
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
	msgChan  chan string
	Trainers [2]Trainer
	Winner   *Trainer
}

func NewBattle(trainer1, trainer2 Trainer, msgChan chan string) *Battle {
	return &Battle{
		Trainers: [2]Trainer{trainer1, trainer2},
		msgChan:  msgChan,
	}
}

func (b *Battle) Run() {
	defer close(b.msgChan)

	trainer1 := &b.Trainers[0]
	trainer2 := &b.Trainers[1]

	pokemon1 := &trainer1.Pokemon[0]
	pokemon2 := &trainer2.Pokemon[0]

	var pokemon1AttackTimeout, pokemon2AttackTimeout float64 = 5, 5
	if pokemon1.Stats.Speed > pokemon2.Stats.Speed {
		pokemon1AttackTimeout = float64(pokemon2.Stats.Speed) / float64(pokemon1.Stats.Speed) * 5
		if pokemon1AttackTimeout < 2 {
			pokemon1AttackTimeout = 2
		}
	} else {
		pokemon2AttackTimeout = float64(pokemon1.Stats.Speed) / float64(pokemon2.Stats.Speed) * 5
		if pokemon2AttackTimeout < 2 {
			pokemon2AttackTimeout = 2
		}
	}

	pokemon1Ticker := time.NewTicker(time.Duration(pokemon1AttackTimeout*1000) * time.Millisecond)
	pokemon2Ticker := time.NewTicker(time.Duration(pokemon2AttackTimeout*1000) * time.Millisecond)

	defer pokemon1Ticker.Stop()
	defer pokemon2Ticker.Stop()

	for {
		select {
		case <-pokemon1Ticker.C:
			msg, err := attack(pokemon1, pokemon2)
			if err != nil {
				pokemon1.Stats.Hp = 0
			}
			b.msgChan <- "\033[32m" + msg
		case <-pokemon2Ticker.C:
			msg, err := attack(pokemon2, pokemon1)
			if err != nil {
				pokemon2.Stats.Hp = 0
			}
			b.msgChan <- "\033[31m" + msg
		}

		if pokemon1.Stats.Hp <= 0 {
			b.msgChan <- trainer1.Name + "'s " + pokemon1.Name + " fainted!"
			b.Winner = trainer2
			break
		}

		if pokemon2.Stats.Hp <= 0 {
			b.msgChan <- trainer2.Name + "'s " + pokemon2.Name + " fainted!"
			b.Winner = trainer1
			break
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

func attack(attacker, defender *pokeutils.Pokemon) (string, error) {
	move, ok := selectMove(attacker)
	if !ok {
		return "", fmt.Errorf("%s has no moves left!", attacker.Name)
	}

	if move.Accuracy < rand.Intn(100) {
		return fmt.Sprintf("%s used %s but missed!\n", attacker.Name, move.Name), nil
	}

	damage, flavourText := pokeutils.CalculateDamage(*attacker, *defender, move)
	msg := fmt.Sprintf("%s used %s and dealt %d damage to %s\n", attacker.Name, move.Name, damage, defender.Name)
	if flavourText != "" {
		msg += flavourText + "\n"
	}

	defender.Stats.Hp -= damage
	return msg, nil
}
