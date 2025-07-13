package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	errExit  = errors.New("exit")
	commands map[string]cliCommand
	pokedex  = make(map[string]Pokemon)
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args ...string) error
}

type config struct {
	next string
	prev string
}

func main() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Help the Pokedex Commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show Pokemon by location area with paginate",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show back Pokemon by location area with paginate",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "It takes the name of a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "It takes the name of a Pokemon",
			callback:    commandCatchPokemon,
		},
		"inspect": {
			name:        "inspect",
			description: "Show information from catch pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Menampilkan semua pokemon yg di dapatkan",
			callback:    commandPokedex,
		},
	}

	scanData := bufio.NewScanner(os.Stdin)
	cfg := new(config)

	for {
		fmt.Print("Pokedex > ")
		if !scanData.Scan() {
			break
		}

		//get input from user
		text := scanData.Text()
		// trim space from getting input user
		cleaned := strings.TrimSpace(text)
		// splits command and args from input user
		fields := strings.Fields(cleaned)

		if len(fields) == 0 {
			continue
		}

		cmdName := strings.ToLower(fields[0])
		args := fields[1:]

		cliCmd, ok := commands[cmdName]
		if !ok {
			fmt.Println("Unknown Coammnd!!, ", cmdName)
			continue
		}

		// run command callback.
		err := cliCmd.callback(cfg, args...)

		switch {
		case err == nil:
			continue
		case errors.Is(err, errExit):
			fmt.Fprintln(os.Stdout, "Closing the Pokedex... Goodbye!")
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Error executing command '%s': %v\n", cmdName, err)
		}
	}

}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for key := range pokedex {
		fmt.Printf("- %s\n", key)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {

	val, ok := pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", args[0])
	fmt.Printf("Height: %v\n", val.Height)
	fmt.Printf("Weight: %v\n", val.Weight)
	fmt.Println("Stats:")
	for _, stat := range val.Stats {
		fmt.Printf("-%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range val.Types {
		fmt.Printf("- %s\n", t.Type.Name)

	}

	return nil
}

func commandCatchPokemon(cfg *config, args ...string) error {

	if len(args) < 1 {
		return errors.New("missing location-area argument")
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", args[0])

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var result Pokemon

	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Unmarshal Error:", err)
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	time.Sleep(500 * time.Millisecond)

	if tryCatchPokemon(result.BaseExperience) {
		fmt.Printf("%s was caught!\n", args[0])
		pokedex[args[0]] = result
	} else {
		fmt.Printf("%s escaped!!\n", args[0])
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {

	if len(args) < 1 {
		return errors.New("missing location-area argument")
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", args[0])

	if cacheData, found := cacheEntry.Get(url); found {
		var data EncounterResponse

		err := json.Unmarshal(cacheData, &data)
		if err != nil {
			return err
		}

		for _, pok := range data.PokemonEncounters {
			fmt.Printf("- %s\n", pok.Pokemon.Name)
		}
		return nil

	}

	res, err := http.Get(url)
	if err != nil {
		return errors.New("failed getting data")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("failed to read body")
	}

	cacheEntry.Set(url, body)

	var data EncounterResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.New("failed unmarshal dataaa")
	}

	for _, pok := range data.PokemonEncounters {
		fmt.Printf("- %s\n", pok.Pokemon.Name)
	}

	return nil

}

func commandMapb(cfg *config, args ...string) error {
	url := cfg.prev

	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	data, err := getPokeApi(url)
	if err != nil {
		return err
	}

	for _, pArea := range data.Results {
		fmt.Println(pArea.Name)
	}

	cfg.next = data.Next
	cfg.prev = data.Previous

	return nil
}

func commandMap(cfg *config, args ...string) error {

	url := cfg.next

	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	data, err := getPokeApi(url)
	if err != nil {
		return err
	}

	for _, pArea := range data.Results {
		fmt.Println(pArea.Name)
	}

	cfg.next = data.Next
	cfg.prev = data.Previous

	return nil

}

func commandExit(cfg *config, args ...string) error {
	return errExit

}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Print("Usage: \n\n")
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return nil
}
