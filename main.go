package main

import (
	"bufio"
	"fmt"
	"github.com/cc-jose-nieto/go-pokedex/internal/PokeApi"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Config struct {
	PokeApiUrl string
	Next       string
	Previous   string
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var actions = map[string]cliCommand{}

func main() {
	godotenv.Load()
	c := Config{}
	c.PokeApiUrl = os.Getenv("POKEAPI_URL")
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
	url := fmt.Sprintf("%s/location-area?offset=0&limit=20", c.PokeApiUrl)

	if c.Next != "" {
		url = c.Next
	}

	res, err := PokeApi.GetLocations(url)

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

	if c.Previous == "" {
		return fmt.Errorf("no previous location found")
	}

	res, err := PokeApi.GetLocations(c.Previous)

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
