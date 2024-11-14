package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/grvbrk/pokedexcli/cache"
)

type CliCommands struct {
	Name     string
	Desc     string
	Callback func(cache *cache.Cache, args ...string)
}

type config struct {
	Offset int
	Limit  int
}

type locationResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokemonsFoundAtLocation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

var store = config{
	Offset: 0,
	Limit:  20,
}

func GetCommands() map[string]CliCommands {

	return map[string]CliCommands{
		"help": {
			Name:     "help",
			Desc:     "Displays a help message",
			Callback: CommandHelp,
		},
		"exit": {
			Name:     "exit",
			Desc:     "Exit the Pokedex",
			Callback: CommandExit,
		},
		"map": {
			Name:     "next",
			Desc:     "Show next map locations",
			Callback: CommandMap,
		},
		"mapb": {
			Name:     "previous",
			Desc:     "Show previous map locations",
			Callback: CommandMapb,
		},
		"explore": {
			Name:     "explore",
			Desc:     "Show all Pokémons in a given area.",
			Callback: CommandExplore,
		},
	}
}

func CommandHelp(cache *cache.Cache, args ...string) {
	fmt.Println("\nWelcome to the Pokedex!\nHere are the available commands:\n1. help\n2. exit\n3. map\n4. mapb\n5. explore")
}

func CommandExit(cache *cache.Cache, args ...string) {
	fmt.Println(`Exiting...`)
	os.Exit(0)
}

func CommandMap(cache *cache.Cache, args ...string) {
	fmt.Println(`Here are the 20 new locations on the world map`)
	baseUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%v&limit=%v", store.Offset, store.Limit)
	val, ok := cache.Get(baseUrl)
	if ok {
		var locations locationResponse
		err := json.Unmarshal(val, &locations)
		if err != nil {
			fmt.Println("Something went wrong...")
		}
		renderLocations(locations)
		return
	}
	response, err := http.Get(baseUrl)
	if err != nil {
		fmt.Println("Couldn't fetch the locations.")
	}
	store.Offset += 20
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Something went wrong...")
	}
	var locations locationResponse
	err = json.Unmarshal(data, &locations)
	if err != nil {
		fmt.Println("Something went wrong...")
	}
	renderLocations(locations)
}

func CommandMapb(cache *cache.Cache, args ...string) {
	store.Offset -= 40
	if store.Offset < 0 {
		store.Offset = 0
		fmt.Println("There's nothing to show. Use 'map' command instead.")
		return
	}
	baseUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location/?offset=%v&limit=%v", store.Offset, store.Limit)
	val, ok := cache.Get(baseUrl)
	if ok {
		var locations locationResponse
		err := json.Unmarshal(val, &locations)
		if err != nil {
			fmt.Println("Something went wrong...")
		}
		renderLocations(locations)
		return
	}
	response, err := http.Get(baseUrl)
	if err != nil {
		fmt.Println("Couldn't fetch the locations.")
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Something went wrong...")
	}
	var locations locationResponse
	err = json.Unmarshal(data, &locations)
	if err != nil {
		fmt.Println("Something went wrong...")
	}
	renderLocations(locations)
}

func renderLocations(locations locationResponse) {
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
}

func CommandExplore(cache *cache.Cache, args ...string) {
	fmt.Printf("Exploring %v...\n", args[0])
	baseUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", args[0])
	response, err := http.Get(baseUrl)
	if err != nil {
		fmt.Println("Couldn't find any pokemon at the location.")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Couldn't find any Pokémon at the location.")
		return
	}

	var pokemons pokemonsFoundAtLocation
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&pokemons)
	if err != nil {
		fmt.Println("Something went wrong...", err)
	}
	RenderPokemons(pokemons)

}

func RenderPokemons(locations pokemonsFoundAtLocation) {
	for _, location := range locations.PokemonEncounters {
		fmt.Println(location.Pokemon.Name)
	}
}
