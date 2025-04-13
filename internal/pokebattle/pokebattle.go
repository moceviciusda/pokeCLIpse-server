package pokebattle

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/moceviciusda/pokeCLIpse-server/pkg/ansiiutils"
	"github.com/moceviciusda/pokeCLIpse-server/pkg/pokeutils"
)

type Pokemon struct {
	pokeutils.Pokemon
	ExpGain int
	BaseExp int
}

type Trainer struct {
	Name    string
	Pokemon []Pokemon

	activePokemon *Pokemon
	participants  []*Pokemon
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
	trainer1.activePokemon = &trainer1.Pokemon[0]
	trainer1.participants = []*Pokemon{&trainer1.Pokemon[0]}
	trainer2.activePokemon = &trainer2.Pokemon[0]
	trainer2.participants = []*Pokemon{&trainer2.Pokemon[0]}
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
	if b.Trainers[0].activePokemon.Stats.Speed > b.Trainers[1].activePokemon.Stats.Speed {
		pokemon1AttackTimeout = float64(b.Trainers[1].activePokemon.Stats.Speed) / float64(b.Trainers[0].activePokemon.Stats.Speed) * 5
		if pokemon1AttackTimeout < 2 {
			pokemon1AttackTimeout = 2
		}
	} else {
		pokemon2AttackTimeout = float64(b.Trainers[0].activePokemon.Stats.Speed) / float64(b.Trainers[1].activePokemon.Stats.Speed) * 5
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
		b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: "A wild " + trainer2.activePokemon.Name + " appeared!"}
	} else {
		b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer2.Name + " sent out " + trainer2.activePokemon.Name + "!"}
	}
	b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer1.Name + " sent out " + trainer1.activePokemon.Name + "!"}

	b.calculateTickers()

	for {
		select {
		case <-b.pokemon1Ticker.C:
			msg, err := attack(trainer1.activePokemon, trainer2.activePokemon)
			if err != nil {
				trainer1.activePokemon.Stats.Hp = 0
				msg = err.Error()
			}
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: msg, Subject: trainer1.Name}
		case <-b.pokemon2Ticker.C:
			msg, err := attack(trainer2.activePokemon, trainer1.activePokemon)
			if err != nil {
				trainer2.activePokemon.Stats.Hp = 0
				msg = err.Error()
			}
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: msg, Subject: trainer2.Name}
		}

		if trainer1.activePokemon.Stats.Hp <= 0 {
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer1.Name + "'s " + trainer1.activePokemon.Name + " fainted!"}

			expYield := pokeutils.ExpYield(trainer2.activePokemon.BaseExp, trainer1.activePokemon.Level)
			for _, p := range trainer2.participants {
				p.ExpGain += expYield / len(trainer2.participants)
			}
			trainer2.participants = []*Pokemon{trainer2.activePokemon}

			err := b.SelectPokemon(trainer1)
			if err != nil {
				b.Winner = trainer2
				b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: err.Error()}
				b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer2.Name + " wins!"}
				return
			}
		}

		if trainer2.activePokemon.Stats.Hp <= 0 {
			b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer2.Name + "'s " + trainer2.activePokemon.Name + " fainted!"}

			expYield := pokeutils.ExpYield(trainer2.activePokemon.BaseExp, trainer2.activePokemon.Level)
			for _, p := range trainer1.participants {
				p.ExpGain += expYield / len(trainer1.participants)
			}
			trainer1.participants = []*Pokemon{trainer1.activePokemon}

			err := b.SelectPokemon(trainer2)
			if err != nil {
				b.Winner = trainer1
				b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: err.Error()}
				b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer1.Name + " wins!"}
				return
			}
		}
	}
}

func (b *Battle) SelectPokemon(trainer *Trainer) error {

	livePokemon := trainer.GetLivePokemon()
	// debugging why battle never ends
	fmt.Println("Trainer: ", trainer.Name)
	for _, p := range livePokemon {

		fmt.Println("Live Pokemon: ", p.Name, p.Stats.Hp)
	}
	if len(livePokemon) == 0 {
		return fmt.Errorf("%s has no live pokemon!", trainer.Name)
	}

	var selectedPokemon *Pokemon

	if len(livePokemon) == 1 {
		selectedPokemon = livePokemon[0]
	} else {

		options := make([]string, 0, len(livePokemon))
		for _, p := range livePokemon {
			options = append(options, p.Name)
		}
		b.MsgChan <- BattleMessage{Type: BattleMsgSelect, Subject: trainer.Name, Options: options}

		selected := <-b.MsgChan
		for i, p := range livePokemon {
			if p.Name == selected.Message {
				selectedPokemon = livePokemon[i]

			}
		}
	}
	if selectedPokemon == nil {
		selectedPokemon = livePokemon[0]
	}

	trainer.activePokemon = selectedPokemon
	trainer.participants = []*Pokemon{selectedPokemon}
	b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer.Name + " sent out " + selectedPokemon.Name + "!"}
	b.calculateTickers()

	// TODO: Implement timeout for PVP battles

	// select {
	// case selected := <-b.MsgChan:
	// 	for i, p := range trainer.Pokemon {
	// 		if p.Name == selected.Message {
	// 			return &trainer.Pokemon[i]
	// 		}
	// 	}
	// case <-time.After(10 * time.Second):
	// 	fmt.Println(trainer.Name + " did not select a pokemon in time!")
	// 	b.MsgChan <- BattleMessage{Type: BattleMsgInfo, Message: trainer.Name + " did not select a pokemon in time!"}
	// }

	return nil
}

func (t *Trainer) GetLivePokemon() []*Pokemon {
	livePokemon := make([]*Pokemon, 0, len(t.Pokemon))
	for i := range t.Pokemon {
		if t.Pokemon[i].Stats.Hp > 0 {
			livePokemon = append(livePokemon, &t.Pokemon[i])
		}
	}
	return livePokemon
}

func selectMove(p *Pokemon) (pokeutils.Move, bool) {
	validMoves := make([]pokeutils.Move, 0, len(p.Moves))
	for _, move := range p.Moves {
		if move.PP > 0 {
			validMoves = append(validMoves, move)
		}
	}
	if len(validMoves) == 0 {
		return pokeutils.Move{}, false
	}

	move := validMoves[rand.Intn(len(validMoves))]
	move.PP--
	return move, true
}

func attack(attacker, defender *Pokemon) (string, error) {
	move, ok := selectMove(attacker)
	if !ok {
		return "", fmt.Errorf("%s has no moves left!", attacker.Name)
	}

	if move.Accuracy < rand.Intn(100) {
		return fmt.Sprintf("%s used %s but missed!\n", attacker.Name, move.Name), nil
	}

	damage, flavourText := pokeutils.CalculateDamage(attacker.Pokemon, defender.Pokemon, move)
	msg := fmt.Sprintf("%s used %s and dealt %d damage to %s\n", ansiiutils.StyleBold+attacker.Name+ansiiutils.StyleNoBold, move.Name, damage, ansiiutils.StyleBold+defender.Name+ansiiutils.StyleNoBold)
	if flavourText != "" {
		msg += flavourText + "\n"
	}

	defender.Stats.Hp -= damage
	return msg, nil
}
