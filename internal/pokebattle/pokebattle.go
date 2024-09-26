package pokebattle

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
)

type Trainer struct {
	Name          string
	Pokemon       []pokeutils.Pokemon
	ActivePokemon *pokeutils.Pokemon
}

type PokemonState struct {
	HP, Attack, Defense, SpecialAttack, SpecialDefense, Speed int
}

const (
	BattleMsgInfo = iota
	BattleMsgSelect
	BattleMsgAction
)

type BattleMessage struct {
	Type    int
	Message string
	Subject string
	Options []string
}

type Battle struct {
	MsgChan                        chan BattleMessage
	Trainers                       [2]Trainer
	Winner                         *Trainer
	pokemon1Ticker, pokemon2Ticker *time.Ticker
}

func NewBattle(trainer1, trainer2 Trainer, msgChan chan BattleMessage) *Battle {
	trainer1.ActivePokemon = &trainer1.Pokemon[0]
	trainer2.ActivePokemon = &trainer2.Pokemon[0]
	return &Battle{
		Trainers: [2]Trainer{trainer1, trainer2},
		MsgChan:  msgChan,
	}
}

func (b *Battle) stopTickers() {
	if b.pokemon1Ticker != nil {
		b.pokemon1Ticker.Stop()
	}
	if b.pokemon2Ticker != nil {
		b.pokemon2Ticker.Stop()
	}
}

func (b *Battle) calculateTickers() {
	b.stopTickers()

	var pokemon1AttackTimeout, pokemon2AttackTimeout float64 = 5, 5
	if b.Trainers[0].ActivePokemon.Stats.Speed > b.Trainers[1].ActivePokemon.Stats.Speed {
		pokemon1AttackTimeout = float64(b.Trainers[1].ActivePokemon.Stats.Speed) / float64(b.Trainers[0].ActivePokemon.Stats.Speed) * 5
		if pokemon1AttackTimeout < 2 {
			pokemon1AttackTimeout = 2
		}
	} else {
		pokemon2AttackTimeout = float64(b.Trainers[0].ActivePokemon.Stats.Speed) / float64(b.Trainers[1].ActivePokemon.Stats.Speed) * 5
		if pokemon2AttackTimeout < 2 {
			pokemon2AttackTimeout = 2
		}
	}

	b.pokemon1Ticker = time.NewTicker(time.Duration(pokemon1AttackTimeout*1000) * time.Millisecond)
	b.pokemon2Ticker = time.NewTicker(time.Duration(pokemon2AttackTimeout*1000) * time.Millisecond)
}

func (b *Battle) Run() {
	defer close(b.MsgChan)
	defer b.stopTickers()

	speedChangeChan := make(chan struct{})
	go func() {
		for range speedChangeChan {
			b.calculateTickers()
		}
	}()

	trainer1 := &b.Trainers[0]
	trainer2 := &b.Trainers[1]

	if trainer2.Name == "Wild" {
		b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: "A wild " + trainer2.ActivePokemon.Name + " appeared!"}
	} else {
		b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer2.Name + " sent out " + trainer2.ActivePokemon.Name + "!"}
	}
	b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer1.Name + " sent out " + trainer1.ActivePokemon.Name + "!"}

	b.calculateTickers()

	for {
		select {
		case <-b.pokemon1Ticker.C:
			msg, err := attack(trainer1.ActivePokemon, trainer2.ActivePokemon)
			if err != nil {
				trainer1.ActivePokemon.Stats.Hp = 0
			}
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: msg, Subject: trainer1.Name}
		case <-b.pokemon2Ticker.C:
			msg, err := attack(trainer2.ActivePokemon, trainer1.ActivePokemon)
			if err != nil {
				trainer2.ActivePokemon.Stats.Hp = 0
			}
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: msg, Subject: trainer2.Name}
		}

		if trainer1.ActivePokemon.Stats.Hp <= 0 {
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer1.Name + "'s " + trainer1.ActivePokemon.Name + " fainted!"}
			pokemon := trainer1.GetLivePokemon()
			if len(pokemon) == 0 {
				b.Winner = trainer2
				b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer2.Name + " wins!"}
				return
			}

			trainer1.ActivePokemon = b.SelectPokemon(*trainer1)
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer1.Name + " sent out " + trainer1.ActivePokemon.Name + "!"}
			b.calculateTickers()
		}

		if trainer2.ActivePokemon.Stats.Hp <= 0 {
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer2.Name + "'s " + trainer2.ActivePokemon.Name + " fainted!"}
			pokemon := trainer2.GetLivePokemon()
			if len(pokemon) == 0 {
				b.Winner = trainer1
				b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer1.Name + " wins!"}
				return
			}

			trainer2.ActivePokemon = b.SelectPokemon(*trainer2)
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer2.Name + " sent out " + trainer2.ActivePokemon.Name + "!"}
			b.calculateTickers()
		}

	}
}

func (b *Battle) SelectPokemon(trainer Trainer) *pokeutils.Pokemon {
	options := make([]string, 0, len(trainer.Pokemon))
	for _, p := range trainer.GetLivePokemon() {
		options = append(options, p.Name)
	}
	b.MsgChan <- BattleMessage{Type: BattleMsgSelect, Subject: trainer.Name, Options: options}

	selected := <-b.MsgChan
	for i, p := range trainer.Pokemon {
		if p.Name == selected.Message {
			return &trainer.Pokemon[i]
		}
	}
	// select {
	// case selected := <-b.MsgChan:
	// 	for i, p := range trainer.Pokemon {
	// 		if p.Name == selected.Message {
	// 			return &trainer.Pokemon[i]
	// 		}
	// 	}
	// case <-time.After(10 * time.Second):
	// 	b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer.Name + " did not select a pokemon in time!"}
	// }

	return &trainer.Pokemon[0]
}

func (t *Trainer) GetLivePokemon() []pokeutils.Pokemon {
	livePokemon := make([]pokeutils.Pokemon, 0, len(t.Pokemon))
	for _, p := range t.Pokemon {
		if p.Stats.Hp > 0 {
			livePokemon = append(livePokemon, p)
		}
	}
	return livePokemon
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
