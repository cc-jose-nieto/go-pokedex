package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var actions = map[string]cliCommand{}

func main() {
	actions = map[string]cliCommand{
		"exit": {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"help": {name: "help", description: "Show available commands", callback: helpCommand},
	}
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, action := range actions {
		fmt.Printf("%s: %s\n", action.name, action.description)
	}
	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	for {

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

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func helpCommand() error {
	fmt.Print("Available commands:\n\n")
	for _, action := range actions {
		fmt.Printf("%s: %s\n", action.name, action.description)
	}
	return nil
}
