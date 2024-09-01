package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (c *Client) GetPokemon(name string) (PokemonResponse, error) {
	url := baseURL + "/pokemon/" + strings.ToLower(name)
	body, ok := c.cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return PokemonResponse{}, err
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return PokemonResponse{}, fmt.Errorf("response failed with status code: %v and body: %s", res.StatusCode, body)
		}
		if err != nil {
			return PokemonResponse{}, err
		}

		c.cache.Add(url, body)
	}

	pokemon := PokemonResponse{}
	err := json.Unmarshal(body, &pokemon)
	if err != nil {
		return PokemonResponse{}, err
	}

	return pokemon, nil
}
