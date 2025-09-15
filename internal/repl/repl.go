package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/evanwiseman/pokedexcli/internal/pokeapi"
)

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func CommandExit(ctx CommandContext) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("error exiting program")
}

func CommandHelp(ctx CommandContext) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, v := range GetCommandRegistry() {
		fmt.Printf("%v: %v\n", v.Name, v.Description)
	}

	return nil
}

func CommandMap(ctx CommandContext) error {
	if ctx.LocationConfig.Next == nil {
		return fmt.Errorf("you're on the last page")
	}
	url := *ctx.LocationConfig.Next
	areas, err := pokeapi.GetLocationAreas(url)
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

func CommandMapb(ctx CommandContext) error {
	if ctx.LocationConfig.Previous == nil {
		return fmt.Errorf("you're on the first page")
	}
	url := *ctx.LocationConfig.Previous
	areas, err := pokeapi.GetLocationAreas(url)
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

type CliCommand struct {
	Name        string
	Description string
	Callback    func(ctx CommandContext) error
}

type CommandContext struct {
	LocationConfig *pokeapi.Config
}

func strPtr(s string) *string {
	return &s
}

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
	}
}

func Start() {
	userInputScanner := bufio.NewScanner(os.Stdin)
	ctx := CommandContext{
		LocationConfig: &pokeapi.Config{
			Next:     strPtr(pokeapi.LocationAreasURL),
			Previous: nil,
		},
	}
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

		// Get command and perform command from registry
		command := tokens[0]
		cli, ok := GetCommandRegistry()[command]
		if !ok {
			fmt.Printf("error command '%s' not in registry\n", command)
			continue
		}
		if err := cli.Callback(ctx); err != nil {
			fmt.Printf("error command failed %v\n", err)
		}
	}
}
