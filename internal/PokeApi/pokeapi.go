package PokeApi

import (
	"encoding/json"
	"fmt"
	"github.com/cc-jose-nieto/go-pokedex/internal/pokecache"
	"io"
	"net/http"
)

type location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
}

func GetLocations(apiUrl string, cache *pokecache.Cache) (LocationResponse, error) {
	var body []byte

	if cached, ok := cache.Get(apiUrl); ok {
		body = cached
	} else {
		res, _ := http.Get(apiUrl)
		bodyIo, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("error reading response body")
			return LocationResponse{}, err
		}

		defer res.Body.Close()

		body = bodyIo

		cache.Add(apiUrl, body)
	}

	var locations LocationResponse

	err := json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println("error unmarshalling response body")
		return LocationResponse{}, err
	}

	return locations, nil
}
