package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var registry map[string]cliCommand

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type area struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type areaResponse struct {
	Results  []area `json:"results"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	return strings.Fields(lower)
}

func commandExit(conf *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, command := range registry {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(conf *config) error {
	res, err := http.Get(conf.Next)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	areas := areaResponse{}
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

func commandMapBack(conf *config) error {
	if conf.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := http.Get(conf.Previous)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	areas := areaResponse{}
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
