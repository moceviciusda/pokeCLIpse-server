package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (c *Client) GetPokemonSpecies(nameOrId string, url string) (PokemonSpeciesResponse, error) {
	if url == "" {
		url = baseURL + "/pokemon-species/" + strings.ToLower(nameOrId)
	}

	body, ok := c.cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return PokemonSpeciesResponse{}, err
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return PokemonSpeciesResponse{}, fmt.Errorf("response failed with status code: %v and body: %s", res.StatusCode, body)
		}
		if err != nil {
			return PokemonSpeciesResponse{}, err
		}

		c.cache.Add(url, body)
	}

	pokemonSpecies := PokemonSpeciesResponse{}
	err := json.Unmarshal(body, &pokemonSpecies)
	if err != nil {
		return PokemonSpeciesResponse{}, err
	}

	return pokemonSpecies, nil
}
