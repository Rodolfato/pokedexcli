package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"

	"github.com/rodolfato/pokedexcli/pokecache"
)

type command struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	next            *string
	previous        *string
	locationsCache  *pokecache.Cache
	encountersCache *pokecache.Cache
	pokemonCaught   map[string]pokemon
}

type resBody struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
}

type locationArea struct {
	Name              string
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

type pokemonEncounter struct {
	Pokemon struct {
		Name string
		Url  string
	}
}

type pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func getCommands() map[string]command {
	return map[string]command{
		"help": {
			name:        "help",
			description: "Displays this message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Gives the next list of in game locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "map",
			description: "Gives the previous list of in game locations",
			callback:    commandMapB,
		},

		"explore": {
			name:        "explore",
			description: "Explores a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func (config *config) modifyConfig(resBody resBody, body []byte, entryName string) {
	config.locationsCache.Add(entryName, body)
	*config.next = resBody.Next
	*config.previous = resBody.Previous
}

func getRequest(url string) ([]byte, error) {
	var zeroValue []byte
	response, err := http.Get(url)
	if err != nil {
		return zeroValue, err
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode > 299 {
		return zeroValue, errors.New("failed request")
	}
	if err != nil {
		return zeroValue, errors.New("error reading request body")
	}
	return body, nil
}

func commandPokedex(config *config, paramters []string) error {
	printCurrentPokemon(config.pokemonCaught)
	return nil
}

func commandInspect(config *config, parameters []string) error {
	pokemon, ok := config.pokemonCaught[parameters[0]]
	if ok {
		fmt.Println()
		fmt.Printf("-name: %s\n", pokemon.Name)
		fmt.Printf("-height: %d\n", pokemon.Height)
		fmt.Printf("-weight: %d\n", pokemon.Weight)
		fmt.Println("stats:")
		for _, entry := range pokemon.Stats {
			fmt.Printf("\t-%s: %d\n", entry.Stat.Name, entry.BaseStat)
		}
		fmt.Println("types:")
		for _, entry := range pokemon.Types {
			fmt.Printf("\t-%s\n", entry.Type.Name)
		}
		return nil
	}
	fmt.Printf("You haven't caught a %s\n", parameters[0])
	return nil
}

func commandMap(config *config, parameters []string) error {

	if *config.next == "" {
		return errors.New("next route non existant")
	}
	fmt.Printf("\tLa siguiente ruta a acceder es %s\n", *config.next)

	body, ok := config.locationsCache.Get(*config.next)
	if !ok {
		var err error
		body, err = getRequest(*config.next)
		if err != nil {
			return err
		}
	}

	resBody := resBody{}

	error_marsh := json.Unmarshal(body, &resBody)
	if error_marsh != nil {
		return error_marsh
	}
	printLocations(resBody)
	fmt.Printf("\tSe modificara el config.next con el siguiente dato %s\n", resBody.Next)
	fmt.Printf("\tEl presente request se obtuvo entrando con la siguiente ruta: %s\n", *config.next)
	config.modifyConfig(resBody, body, *config.next)
	return nil
}

func commandMapB(config *config, parameters []string) error {
	if *config.previous == "" {
		return errors.New("previous route non existant")
	}
	fmt.Printf("\tLa anterior ruta a acceder es %s\n", *config.previous)
	body, ok := config.locationsCache.Get(*config.previous)
	if !ok {
		var err error
		body, err = getRequest(*config.previous)
		if err != nil {
			return err
		}
	}
	resBody := resBody{}

	error_marsh := json.Unmarshal(body, &resBody)
	if error_marsh != nil {
		return error_marsh
	}
	printLocations(resBody)
	fmt.Printf("\tSe modificara el config.previous con el siguiente dato %s\n", resBody.Previous)
	fmt.Printf("\tEl presente request se obtuvo entrando con la siguiente ruta: %s\n", *config.previous)
	config.modifyConfig(resBody, body, *config.previous)
	return nil
}

func printEncounters(locationAreaEncounters locationArea) {
	for _, result := range locationAreaEncounters.PokemonEncounters {
		fmt.Println(result.Pokemon.Name)
	}
}

func printCurrentPokemon(pokemonCaught map[string]pokemon) {
	fmt.Println("\nThese are my pokemon!")
	for _, entry := range pokemonCaught {
		fmt.Println(entry.Name)
	}
}

func tryToCatch(config *config, pokemon pokemon) {
	fmt.Printf("%s found!\n", pokemon.Name)
	baseExperience := pokemon.BaseExperience
	tries := 3
	for baseExperience > 0 && tries > 0 {
		random := rand.IntN(100)
		fmt.Printf("Attempt n° %d. Base experience: %d. Random %d.\n", tries, baseExperience, random)
		baseExperience -= random
		tries -= 1
	}
	if baseExperience <= 0 {
		fmt.Printf("%d base experience\n", baseExperience)
		fmt.Printf("%s caught!\n", pokemon.Name)
		config.pokemonCaught[pokemon.Name] = pokemon
		printCurrentPokemon(config.pokemonCaught)
	} else if tries <= 0 {
		fmt.Printf("%d base experience\n", baseExperience)
		fmt.Printf("%s escaped!\n", pokemon.Name)
		printCurrentPokemon(config.pokemonCaught)
	}

}

func commandCatch(config *config, parameters []string) error {
	fmt.Println("I have entered catch and these are the parameters I've got: ")
	for i, parameter := range parameters {
		fmt.Printf("The %d parameter is %s\n", i, parameter)
	}

	apiCall := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", parameters[0])

	body, err := getRequest(apiCall)
	if err != nil {
		return err
	}

	pokemon := pokemon{}

	error_marsh := json.Unmarshal(body, &pokemon)
	if error_marsh != nil {
		return error_marsh
	}

	tryToCatch(config, pokemon)
	return nil
}

func commandExplore(config *config, parameters []string) error {
	// API call https://pokeapi.co/api/v2/location-area/{id or name}/
	fmt.Println("I have entered explore and these are the parameters I've got: ")
	for i, parameter := range parameters {
		fmt.Printf("The %d parameter is %s\n", i, parameter)
	}

	apiCall := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", parameters[0])
	body, ok := config.encountersCache.Get(parameters[0])
	if !ok {
		var err error
		body, err = getRequest(apiCall)
		config.encountersCache.Add(parameters[0], body)
		if err != nil {
			return err
		}
	}

	locationAreaBody := locationArea{}

	error_marsh := json.Unmarshal(body, &locationAreaBody)
	if error_marsh != nil {
		return error_marsh
	}
	printEncounters(locationAreaBody)
	return nil
}

func printLocations(resBody resBody) {
	for _, result := range resBody.Results {
		fmt.Println(result.Name)
	}
}

func getUserInput() (string, []string) {

	reader := bufio.NewReader(os.Stdin)
	printPrompt()
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	parts := strings.Split(input, " ")
	return parts[0], parts[1:]
}

func printPrompt() {
	fmt.Print("Pokedex > ")
}

func printCmdNotFound(cmdName string) {
	fmt.Printf("Command %s not found\n", cmdName)
}

func commandHelp(config *config, parameters []string) error {

	for _, element := range getCommands() {
		fmt.Printf("%s: %s\n", element.name, element.description)
	}
	return nil
}

func commandExit(config *config, parameters []string) error {
	os.Exit(0)
	return nil
}

func initConfig() *config {
	url := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	prev := ""
	pokemonCaught := make(map[string]pokemon)
	config := config{
		next:            &url,
		previous:        &prev,
		locationsCache:  pokecache.NewCache(10 * pokecache.Second),
		encountersCache: pokecache.NewCache(30 * pokecache.Second),
		pokemonCaught:   pokemonCaught,
	}

	return &config
}

func main() {
	config := initConfig()
	pokecache.Hello()
	for {
		cmdName, parameters := getUserInput()
		cmdFound, ok := getCommands()[cmdName]

		fmt.Printf("The command name is %s\n", cmdFound.name)
		for i, parameter := range parameters {
			fmt.Printf("The parameter %d is %s\n", i, parameter)
		}

		if ok {
			err := cmdFound.callback(config, parameters)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			printCmdNotFound(cmdName)
		}
	}

}
