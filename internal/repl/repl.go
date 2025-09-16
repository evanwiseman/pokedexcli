package repl

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/evanwiseman/pokedexcli/internal/pokeapi"
)

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

type Context struct {
	Client         *pokeapi.Client
	LocationConfig *pokeapi.Config
	Pokedex        map[string]pokeapi.Pokemon
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(ctx *Context, parameters []string) error
}

// Registry containing all repl commands, map of command -> name, description, callback
func GetCommandRegistry() map[string]CliCommand {
	return map[string]CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    CommandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Gets next 20 map locations",
			Callback:    CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Gets previous 20 map locations",
			Callback:    CommandMapb,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore an area, gets a list of all pokemon in the area",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Try to catch a Pokemon.",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a caught Pokemon",
			Callback:    CommandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Lists all caught Pokemon in your Pokedex",
			Callback:    CommandPokedex,
		},
	}
}

// Exits the program
func CommandExit(ctx *Context, parameters []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("error exiting program")
}

// Outputs registry commands and their descriptions
func CommandHelp(ctx *Context, parameters []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, v := range GetCommandRegistry() {
		fmt.Printf("%v: %v\n", v.Name, v.Description)
	}

	return nil
}

// Gets the next 20 areas from location-areas
func CommandMap(ctx *Context, parameters []string) error {
	if ctx.LocationConfig.Next == nil {
		return fmt.Errorf("you're on the last page")
	}
	url := *ctx.LocationConfig.Next
	areas, err := ctx.Client.GetLocationAreaList(url)
	if err != nil {
		return err
	}

	ctx.LocationConfig.Next = areas.Next
	ctx.LocationConfig.Previous = areas.Previous

	for _, result := range areas.Results {
		fmt.Printf("%s\n", result.Name)
	}

	return nil
}

// Gets the previous 20 areas from location-areas
func CommandMapb(ctx *Context, parameters []string) error {
	if ctx.LocationConfig.Previous == nil {
		return fmt.Errorf("you're on the first page")
	}
	url := *ctx.LocationConfig.Previous
	areas, err := ctx.Client.GetLocationAreaList(url)
	if err != nil {
		return err
	}

	ctx.LocationConfig.Next = areas.Next
	ctx.LocationConfig.Previous = areas.Previous

	for _, result := range areas.Results {
		fmt.Printf("%s\n", result.Name)
	}

	return nil
}

// Explores a provided area and lists Pokemon in the area
func CommandExplore(ctx *Context, parameters []string) error {
	if len(parameters) == 0 {
		return fmt.Errorf("'explore' no area provided")
	}
	if len(parameters) > 1 {
		return fmt.Errorf("'explore' expects only one area. try replacing ' ' with '-'")
	}
	area, err := ctx.Client.GetLocationArea(parameters[0])
	if err != nil {
		return err
	}
	for _, encounter := range area.PokemonEncounters {
		fmt.Printf("%s\n", encounter.Pokemon.Name)
	}

	return nil
}

// Attempts to catch a Pokemon and add to users Pokedex
func CommandCatch(ctx *Context, parameters []string) error {
	if len(parameters) == 0 {
		return fmt.Errorf("'catch' no pokemon provided")
	}
	if len(parameters) > 1 {
		return fmt.Errorf("'catch' expects only one pokemon. try replacing ' ' with '-'")
	}

	// Grab pokemon from the Pokedex
	key := parameters[0]
	pokemon, err := ctx.Client.GetPokemon(key)
	if err != nil {
		return err
	}

	const maxBaseExp = 635

	// calculate catch probability clamp between 0.1 and 0.9
	prob := 1.0 - float64(pokemon.BaseExperience)/float64(maxBaseExp) // higher baseExp = lower prob
	if prob < 0.1 {
		prob = 0.1
	}
	if prob > 0.9 {
		prob = 0.9
	}
	// add to pokedex if caught, otherwise let use know it failed
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)
	if rand.Float64() < prob { // Success
		ctx.Pokedex[key] = *pokemon
		fmt.Printf("%v was caught!\n", pokemon.Name)
		fmt.Println("You may now inspect it with the inspect command.")
	} else { // Failure
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}
	return nil
}

// Inspect properties of a Pokemon in the users Pokedex
func CommandInspect(ctx *Context, parameters []string) error {
	if len(parameters) == 0 {
		return fmt.Errorf("'inspect' no pokemon provided")
	}
	if len(parameters) > 1 {
		return fmt.Errorf("'inspect' expects only one pokemon. try replacing ' ' with '-'")
	}

	// Grab pokemon from the Pokedex
	key := parameters[0]
	pokemon, ok := ctx.Pokedex[key]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}

	// Output pertinent information about the Pokemon
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, item := range pokemon.Stats {
		fmt.Printf("  - %v: %v\n", item.Stat.Name, item.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, item := range pokemon.Types {
		fmt.Printf("  - %v\n", item.Type.Name)
	}

	return nil
}

// Lists all Pokemon in the users Pokedex
func CommandPokedex(ctx *Context, parameters []string) error {
	if len(parameters) > 0 {
		return fmt.Errorf("'pokedex' expected no parameters, got %v", parameters)
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range ctx.Pokedex {
		fmt.Printf("  - %v\n", pokemon.Name)
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}

func Start() {
	ctx := Context{
		Client: pokeapi.NewClient(),
		LocationConfig: &pokeapi.Config{
			Next:     strPtr(pokeapi.LocationAreaURL),
			Previous: nil,
		},
		Pokedex: make(map[string]pokeapi.Pokemon),
	}
	userInputScanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		// Block until a user gives input
		success := userInputScanner.Scan()
		if !success {
			continue
		}

		// Get user input and parse into tokens
		text := userInputScanner.Text()
		tokens := CleanInput(text)
		if len(tokens) == 0 {
			continue
		}

		// Get command and perform command from registry
		command := tokens[0]
		cli, ok := GetCommandRegistry()[command]
		if !ok {
			fmt.Printf("error command '%s' not in registry\n", command)
			continue
		}

		// Run the command with the current context
		if err := cli.Callback(&ctx, tokens[1:]); err != nil {
			fmt.Printf("error command failed: %v\n", err)
		}
	}
}
