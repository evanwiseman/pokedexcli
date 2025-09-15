package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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
		tokens := cleanInput(text)

		// Get command and perform command from registry
		command := tokens[0]
		cli, ok := getCommandRegistry()[command]
		if !ok {
			fmt.Printf("error command '%s' not in registry", command)
		}
		if err := cli.callback(); err != nil {
			fmt.Printf("error command failed %v", err)
		}
	}
}
