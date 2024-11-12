package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	r := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		r.Scan()
		words := cleanInput(r.Text())
		if len(words) == 0 {
			continue
		}
		commandName := words[0]
		cmd, ok := getCommands()[commandName]
		if ok {
			cmd.callback()
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}
