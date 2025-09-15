package main

import (
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("error exiting program")
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	help, ok := getCommandRegistry()["help"]
	if !ok {
		return fmt.Errorf("error 'help' not registered")
	}
	fmt.Printf("help: %s\n", help.description)

	exit, ok := getCommandRegistry()["exit"]
	if !ok {
		return fmt.Errorf("error 'exit' not registered")
	}
	fmt.Printf("exit: %s\n", exit.description)
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func getCommandRegistry() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}
