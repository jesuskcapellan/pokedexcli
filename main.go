package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	registry = map[string]cliCommand{
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
			registry[command].callback()
		}
	}
}
