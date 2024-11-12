package main

import (
	"fmt"
	"os"
)

type cliCommands struct {
	name     string
	desc     string
	callback func()
}

func getCommands() map[string]cliCommands {

	return map[string]cliCommands{
		"help": {
			name:     "help",
			desc:     "Displays a help message",
			callback: commandHelp,
		},
		"exit": {
			name:     "exit",
			desc:     "Exit the Pokedex",
			callback: commandExit,
		},
		"map":  {},
		"mapb": {},
	}
}

func commandHelp() {
	fmt.Println("\nWelcome to the Pokedex!\nHere are the available commands:\n1. help\n2. exit")
}

func commandExit() {
	fmt.Println(`Exiting...`)
	os.Exit(0)
}
