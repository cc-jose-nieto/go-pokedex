package PokeApi

import (
	"encoding/json"
	"fmt"
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

func GetLocations(apiUrl string) (LocationResponse, error) {

	res, _ := http.Get(apiUrl)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error reading response body")
		return LocationResponse{}, err
	}

	defer res.Body.Close()

	var locations LocationResponse

	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println("error unmarshalling response body")
		return LocationResponse{}, err
	}

	return locations, nil
}
