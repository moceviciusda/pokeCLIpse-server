package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
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
