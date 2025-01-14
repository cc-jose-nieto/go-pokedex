package main

import (
	"bufio"
	"fmt"
	"github.com/cc-jose-nieto/go-pokedex/internal/PokeApi"
	"github.com/cc-jose-nieto/go-pokedex/internal/PokeBall"
	"github.com/cc-jose-nieto/go-pokedex/internal/Pokedex"
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
	callback    func(config *Config, args ...string) error
}

var actions = map[string]cliCommand{}

var cache *pokecache.Cache = pokecache.NewCache(time.Second * 10)

var pokedex Pokedex.Pokedex = Pokedex.Pokedex{
	Pokemons: make(map[string]PokeApi.Pokemon),
}

func main() {
	godotenv.Load()
	c := Config{}
	c.PokeApiUrl = os.Getenv("POKEAPI_URL")
	c.LocationAreaUrl = fmt.Sprintf("%s/location-area", c.PokeApiUrl)
	c.Next = fmt.Sprintf("%s/location-area", c.PokeApiUrl)
	c.Previous = fmt.Sprintf("%s/location-area", c.PokeApiUrl)

	actions = map[string]cliCommand{
		"exit":    {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"help":    {name: "help", description: "Show available commands", callback: commandHelp},
		"map":     {name: "map", description: "", callback: commandMapLocations},
		"mapb":    {name: "mapb", description: "", callback: commandMapBackLocations},
		"explore": {name: "explore", description: "", callback: commandExplore},
		"catch":   {name: "catch", description: "", callback: commandCatch},
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

		_ = actions[words[0]].callback(&c, words[1:]...)
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

func commandExit(c *Config, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config, args ...string) error {
	fmt.Print("Available commands:\n\n")
	for _, action := range actions {
		fmt.Printf("%s: %s\n", action.name, action.description)
	}
	return nil
}

func commandMapLocations(c *Config, args ...string) error {

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

func commandMapBackLocations(c *Config, args ...string) error {
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

func commandExplore(c *Config, args ...string) error {
	fmt.Printf("Exploring %s...\n", args[0])
	url := fmt.Sprintf("%s/location-area/%s", c.PokeApiUrl, args[0])

	pokemons, err := PokeApi.GetPokemonFromLocationArea(url, cache)
	if err != nil {
		return fmt.Errorf("error getting pokemon: %v", err)
	}

	for _, pokemon := range pokemons {
		fmt.Println(pokemon.Name)
	}

	return nil
}

func commandCatch(c *Config, args ...string) error {
	pokemonName := args[0]

	if pokemonName == "" {
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	if strings.ToLower(pokemonName) == "maluco" {
		time.Sleep(time.Second * 5)
		fmt.Println("You caught a MALUCO!!!")
		return nil
	}

	url := fmt.Sprintf("%s/pokemon/%s", c.PokeApiUrl, pokemonName)

	pokemon, err := PokeApi.GetPokemonByName(url, cache)

	if err != nil {
		fmt.Printf("Pokemon %s does not exist", pokemonName)
		return nil
	}

	time.Sleep(time.Second * 5)
	if ok := PokeBall.Catching(pokemon); ok {
		err = pokedex.Add(pokemon)
		if err != nil {
			fmt.Printf("error adding pokemon to pokedex: %v\n", err)
			return nil
		}
	} else {
		fmt.Println("pokemon not caught, try again")
	}

	fmt.Println(pokemon)

	return nil
}
