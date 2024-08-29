package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
)

func (c *Client) GetLocationAreas(locationUrl string) (LocationAreasResponse, error) {
	var url string
	if locationUrl == "" {
		url = baseURL + "/location-area"
	} else {
		url = locationUrl
	}

	body, ok := c.cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return LocationAreasResponse{}, err
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return LocationAreasResponse{}, fmt.Errorf("response failed with status code: %v and body: %s", res.StatusCode, body)
		}
		if err != nil {
			return LocationAreasResponse{}, err
		}

		c.cache.Add(url, body)
	}

	locations := LocationAreasResponse{}
	err := json.Unmarshal(body, &locations)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	return locations, nil
}

func (c *Client) GetLocationArea(name string) (LocationAreaResponse, error) {
	url := baseURL + "/location-area/" + name

	body, ok := c.cache.Get(url)
	if !ok {
		res, err := c.httpClient.Get(url)
		if err != nil {
			return LocationAreaResponse{}, err
		}

		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			return LocationAreaResponse{}, fmt.Errorf("response failed with status code: %v and body: %s", res.StatusCode, body)
		}
		if err != nil {
			return LocationAreaResponse{}, err
		}

		c.cache.Add(url, body)
	}

	location := LocationAreaResponse{}
	err := json.Unmarshal(body, &location)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	return location, nil
}
