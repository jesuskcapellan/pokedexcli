package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/jesuskcapellan/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	registry = map[string]cliCommand{
		"inspect": {
			name:        "inspect",
			description: "Displays details about a pokemon if it is caught",
			callback:    commandInspect,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon",
			callback:    commandCatch,
		},
		"explore": {
			name:        "explore",
			description: "Displays a list of pokemon in the area",
			callback:    commandExplore,
		},
		"map": {
			name:        "map",
			description: "Displays a list of areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous list of areas",
			callback:    commandMapBack,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	conf := config{
		Next:   "https://pokeapi.co/api/v2/location-area",
		cache:  pokecache.NewCache(5 * time.Second),
		caught: pokecache.NewCache(math.MaxInt), // don't clear cache
	}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			fmt.Println(scanner.Err())
		}
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		command := cleanedInput[0]
		if registry[command].callback == nil {
			fmt.Println("Unknown command")
		} else {
			registry[command].callback(&conf, cleanedInput[1:]...)
		}
	}
}
