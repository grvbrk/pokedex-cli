package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/grvbrk/pokedexcli/cache"
	"github.com/grvbrk/pokedexcli/commands"
)

func main() {
	r := bufio.NewScanner(os.Stdin)
	cache := cache.NewCache(5 * time.Minute)
	for {
		fmt.Print("Pokedex > ")
		r.Scan()
		words := cleanInput(r.Text())
		if len(words) == 0 {
			continue
		}
		switch len(words) {
		case 0:
			continue
		case 1:
			commandName := words[0]
			if commandName == "explore" {
				fmt.Println("Use the map to see locations. Explore locations by passing the name after explore command")
				continue
			}
			cmd, ok := commands.GetCommands()[commandName]
			if ok {
				cmd.Callback(&cache)
				continue
			} else {
				fmt.Println("Unknown command")
				continue
			}
		case 2:
			commandName := words[0]
			arg := words[1]
			cmd, ok := commands.GetCommands()[commandName]
			if ok {
				cmd.Callback(&cache, arg)
				continue
			} else {
				fmt.Println("Unknown command")
				continue
			}
		}

	}
}
