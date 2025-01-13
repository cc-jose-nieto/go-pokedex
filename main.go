package main

import (
	"bufio"
	"fmt"
	"github.com/cc-jose-nieto/go-pokedex/internal/PokeApi"
	"github.com/cc-jose-nieto/go-pokedex/internal/pokecache"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"time"
)

type Config struct {
	PokeApiUrl      string
	LocationAreaUrl string
	Next            string
	Previous        string
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var actions = map[string]cliCommand{}

var cache *pokecache.Cache = pokecache.NewCache(time.Second * 10)

func main() {
	godotenv.Load()
	c := Config{}
	c.PokeApiUrl = os.Getenv("POKEAPI_URL")
	c.LocationAreaUrl = fmt.Sprintf("%s/location-area", c.PokeApiUrl)
	c.Next = fmt.Sprintf("%s/location-area", c.PokeApiUrl)
	c.Previous = fmt.Sprintf("%s/location-area", c.PokeApiUrl)

	actions = map[string]cliCommand{
		"exit": {name: "exit", description: "Exit the Pokedex", callback: func() error { return commandExit(&c) }},
		"help": {name: "help", description: "Show available commands", callback: func() error { return commandHelp(&c) }},
		"map":  {name: "map", description: "", callback: func() error { return commandMapLocations(&c) }},
		"mapb": {name: "mapb", description: "", callback: func() error { return commandMapBackLocations(&c) }},
	}
	//fmt.Print("Welcome to the Pokedex!\n")
	//fmt.Print("Usage:\n\n")
	//for _, action := range actions {
	//	fmt.Printf("%s: %s\n", action.name, action.description)
	//}
	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	for {
		fmt.Print("Pokedex > ")
		input := scanner.Scan()
		if !input {
			fmt.Println("Unknown command")
			continue
		}

		words := cleanInput(scanner.Text())

		if actions[words[0]].name == "" {
			fmt.Println("Unknown command")
			continue
		}

		_ = actions[words[0]].callback()
	}
}

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	words := strings.Split(text, " ")
	var newWords []string
	for i := range words {
		if strings.TrimSpace(words[i]) != "" {
			newWords = append(newWords, words[i])
		}
	}
	return newWords
}

func commandExit(c *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config) error {
	fmt.Print("Available commands:\n\n")
	for _, action := range actions {
		fmt.Printf("%s: %s\n", action.name, action.description)
	}
	return nil
}

func commandMapLocations(c *Config) error {

	res, err := PokeApi.GetLocations(c.Next, cache)

	if err != nil {
		return fmt.Errorf("error getting locations: %v", err)
	}
	c.Next = res.Next
	c.Previous = res.Previous
	for _, location := range res.Results {
		fmt.Println(location.Name)
	}

	return nil

}

func commandMapBackLocations(c *Config) error {
	fmt.Println(c.Previous)
	res, err := PokeApi.GetLocations(c.Previous, cache)

	if err != nil {
		return fmt.Errorf("error getting locations: %v", err)
	}

	c.Next = res.Next
	if c.Previous = res.Previous; c.Previous == "" {
		c.Previous = c.LocationAreaUrl
	}

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}

	return nil

}
