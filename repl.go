package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/jesuskcapellan/pokedexcli/internal/pokecache"
)

var registry map[string]cliCommand

type config struct {
	cache    *pokecache.Cache
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type area struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type areasResponse struct {
	Results  []area `json:"results"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type areaPokemonResponse struct {
	Name string `json:"name"`
}

type areaLocationResponse struct {
	Name string `json:"name"`
}

type areaPokemonEncountersResponse struct {
	Pokemon areaPokemonResponse `json:"pokemon"`
}

type areaResponse struct {
	Name              string                          `json:"name"`
	PokemonEncounters []areaPokemonEncountersResponse `json:"pokemon_encounters"`
	Location          areaLocationResponse            `json:"location"`
}

type pokemonResponse struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	return strings.Fields(lower)
}

func commandExit(conf *config, _ ...string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, _ ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range registry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(conf *config, _ ...string) error {
	var data []byte
	var err error
	if entry, ok := conf.cache.Get(conf.Next); ok {
		data = entry
	} else {
		res, err := http.Get(conf.Next)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		conf.cache.Add(conf.Next, data)
	}
	areas := areasResponse{}
	err = json.Unmarshal(data, &areas)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	conf.Next = areas.Next
	conf.Previous = areas.Previous
	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapBack(conf *config, _ ...string) error {
	if conf.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	var data []byte
	var err error
	if entry, ok := conf.cache.Get(conf.Previous); ok {
		data = entry
	} else {
		res, err := http.Get(conf.Previous)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		conf.cache.Add(conf.Previous, data)
	}
	areas := areasResponse{}
	err = json.Unmarshal(data, &areas)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	conf.Next = areas.Next
	conf.Previous = areas.Previous
	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(conf *config, args ...string) error {
	if len(args) != 1 {
		fmt.Println("Must provide one area name or id")
		return nil
	}
	id := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", id)
	var data []byte
	var err error
	if entry, ok := conf.cache.Get(url); ok {
		data = entry
	} else {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		conf.cache.Add(url, data)
	}
	area := areaResponse{}
	err = json.Unmarshal(data, &area)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("Exploring %s...\n", area.Location.Name)
	fmt.Println("Found Pokemon:")
	for _, area := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", area.Pokemon.Name)
	}
	return nil
}

func commandCatch(conf *config, args ...string) error {
	if len(args) != 1 {
		fmt.Println("Must provide one pokmeon name")
		return nil
	}
	id := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", id)
	var data []byte
	var err error
	if entry, ok := conf.cache.Get(url); ok {
		data = entry
	} else {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		conf.cache.Add(url, data)
	}
	pokemon := pokemonResponse{}
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Printf("Throwing a pokeball at %s...\n", strings.ToTitle(pokemon.Name))
	chance := math.Abs(rand.NormFloat64() * 200)

	if chance >= float64(pokemon.BaseExperience) {
		fmt.Printf("You caught %s!\n", strings.ToTitle(pokemon.Name))
	} else {
		fmt.Printf("%s escaped!\n", strings.ToTitle(pokemon.Name))
	}
	return nil
}
