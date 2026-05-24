package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			fmt.Printf("%s\n", scanner.Err())
		}
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		fmt.Printf("Your command was: %s\n", cleanedInput[0])
	}
}
