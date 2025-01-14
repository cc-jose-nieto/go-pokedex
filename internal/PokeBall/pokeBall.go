package PokeBall

import (
	"github.com/cc-jose-nieto/go-pokedex/internal/PokeApi"
	"math/rand"
)

func Catching(pokemon PokeApi.Pokemon) bool {
	r := rand.Intn(pokemon.BaseExperience * 2)
	if r > pokemon.BaseExperience {
		return true
	}
	return false
}
