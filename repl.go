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
	caught   *pokecache.Cache
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

type pokemonStatResponse struct {
	Name string `json:"name"`
}

type pokemonStatsResponse struct {
	Stat     pokemonStatResponse `json:"stat"`
	BaseStat int                 `json:"base_stat"`
}

type pokemonTypeResponse struct {
	Name string `json:"name"`
}

type pokemonTypesResponse struct {
	Type pokemonTypeResponse `json:"type"`
}

type pokemonResponse struct {
	Name           string                 `json:"name"`
	Height         int                    `json:"height"`
	Weight         int                    `json:"weight"`
	BaseExperience int                    `json:"base_experience"`
	Stats          []pokemonStatsResponse `json:"stats"`
	Types          []pokemonTypesResponse `json:"types"`
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
	data, err := getAndCacheData(conf.Next, conf.cache)
	if err != nil {
		fmt.Println(err.Error())
		return err
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
	data, err := getAndCacheData(conf.Previous, conf.cache)
	if err != nil {
		fmt.Println(err.Error())
		return err
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
	data, err := getAndCacheData(url, conf.cache)
	if err != nil {
		fmt.Println(err.Error())
		return err
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
	data, err := getAndCacheData(url, conf.cache)
	if err != nil {
		fmt.Println(err.Error())
		return err
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
		conf.caught.Add(pokemon.Name, data)
	} else {
		fmt.Printf("%s escaped!\n", strings.ToTitle(pokemon.Name))
	}
	return nil
}

func commandInspect(conf *config, args ...string) error {
	if len(args) != 1 {
		fmt.Println("Must provide one pokmeon name")
		return nil
	}
	name := args[0]
	if entry, ok := conf.caught.Get(name); ok {
		pokemon := pokemonResponse{}
		err := json.Unmarshal(entry, &pokemon)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, pokemonType := range pokemon.Types {
			fmt.Printf("  - %s\n", pokemonType.Type.Name)
		}
		return nil
	}
	fmt.Printf("You have not caught %s.\n", name)
	return nil
}

func commandPokedex(conf *config, _ ...string) error {
	fmt.Println("Your Pokedex:")
	for _, key := range conf.caught.List() {
		if entry, ok := conf.caught.Get(key); ok {
			pokemon := pokemonResponse{}
			err := json.Unmarshal(entry, &pokemon)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			fmt.Printf(" - %s\n", pokemon.Name)
		}
	}
	return nil
}

func getAndCacheData(key string, cache *pokecache.Cache) ([]byte, error) {
	var data []byte
	if entry, ok := cache.Get(key); ok {
		data = entry
	} else {
		res, err := http.Get(key)
		if err != nil {
			return []byte{}, err
		}
		data, err = io.ReadAll(res.Body)
		if err != nil {
			return []byte{}, err
		}
		cache.Add(key, data)
	}
	return data, nil
}
