package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
)

func (c *Client) GetEvolutionChain(url string) (EvolutionChain, error) {

	body, ok := c.cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return EvolutionChain{}, err
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return EvolutionChain{}, fmt.Errorf("response failed with status code: %v and body: %s", res.StatusCode, body)
		}
		if err != nil {
			return EvolutionChain{}, err
		}

		c.cache.Add(url, body)
	}

	evolutionChain := EvolutionChain{}
	err := json.Unmarshal(body, &evolutionChain)
	if err != nil {
		return EvolutionChain{}, err
	}

	return evolutionChain, nil
}
