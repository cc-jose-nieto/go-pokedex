package Pokedex

import (
	"fmt"
	"github.com/cc-jose-nieto/go-pokedex/internal/PokeApi"
)

type Pokedex struct {
	Pokemons map[string]PokeApi.Pokemon
}

func (pokedex *Pokedex) Add(pokemon PokeApi.Pokemon) error {
	pokedex.Pokemons[pokemon.Name] = pokemon
	return nil
}

func (pokedex *Pokedex) Get(pokemonName string) (PokeApi.Pokemon, error) {
	pokemon, ok := pokedex.Pokemons[pokemonName]
	if !ok {
		return PokeApi.Pokemon{}, fmt.Errorf("you have not caught that pokemon")
	}

	return pokemon, nil
}
