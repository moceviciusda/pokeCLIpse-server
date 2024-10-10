package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
)

func (c *Client) GetMove(nameOrId string) (MoveResponse, error) {
	url := baseURL + "/move/" + nameOrId

	body, ok := c.cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return MoveResponse{}, err
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return MoveResponse{}, fmt.Errorf("response failed with status code: %v and body: %s", res.StatusCode, body)
		}
		if err != nil {
			return MoveResponse{}, err
		}

		c.cache.Add(url, body)
	}

	move := MoveResponse{}
	err := json.Unmarshal(body, &move)
	if err != nil {
		return MoveResponse{}, err
	}

	return move, nil
}

func (c *Client) SelectRandomMoves(pokemonNameOrId string, level int) ([]MoveResponse, error) {
	pokemon, err := c.GetPokemon(pokemonNameOrId)
	if err != nil {
		return nil, err
	}

	moveOptions := make(map[string]MoveResponse)
	for _, move := range pokemon.Moves {
		if _, ok := moveOptions[move.Move.Name]; ok {
			continue
		}

		for _, details := range move.VersionGroupDetails {
			if details.LevelLearnedAt > level {
				continue
			}
			if !(details.MoveLearnMethod.Name == "level-up" || details.MoveLearnMethod.Name == "egg") {
				continue
			}

			m, err := c.GetMove(move.Move.Name)
			if err != nil {
				return nil, err
			}

			moveOptions[move.Move.Name] = m
			break
		}
	}

	moves := make([]MoveResponse, 0, 4)
	for i := 0; i < 4; i++ {
		if len(moveOptions) == 0 || len(moves) == 4 {
			break
		}

		moveOptKeys := make([]string, 0, len(moveOptions))
		for k := range moveOptions {
			moveOptKeys = append(moveOptKeys, k)
		}
		moveName := moveOptKeys[rand.Intn(len(moveOptKeys))]

		moves = append(moves, moveOptions[moveName])
		delete(moveOptions, moveName)
	}

	return moves, nil
}

func (c *Client) GetMovesLearnedAtLvl(pokemonNameOrId string, level int) (map[string]MoveResponse, error) {
	pokemon, err := c.GetPokemon(pokemonNameOrId)
	if err != nil {
		return nil, err
	}

	movesLearned := make(map[string]MoveResponse)
	for _, move := range pokemon.Moves {
		for _, details := range move.VersionGroupDetails {
			if details.LevelLearnedAt == level {
				m, err := c.GetMove(move.Move.Name)
				if err != nil {
					return nil, err
				}

				movesLearned[move.Move.Name] = m
				break
			}
		}
	}

	return movesLearned, nil
}
