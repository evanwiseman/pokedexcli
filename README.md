# pokedexcli

This is a pokedex command line that uses the pokeapi to search areas and locations for pokemon, to catch pokemon in the area, and to inspect caught pokemon. Uses Go, http requests, and caching to store responses from the api.

## Requirements
Go v1.24.0+

## Commands
"help" (usage: help) - Displays a help message containing all commands, their description, and their callback
"exit" (usage: exit) - Exits the Pokedex
"map" (usage: map) - Gets the next 20 map locations from the /api/v2/location-area endpoint
"mapb" (usage: mapb) - Gets the previous 20 map locations from the /api/v2/location-area endpoint
"explore" (usage: explore <area>) - Explores the specified area, and lists all pokemon located in the area
"catch" (usage: catch <pokemon>) - Attempts to catch a pokemon located in the area
"inspect" (usage: inspect <pokemon>) - Inspects a caught pokemon in your pokedex
"pokedex" (usage: pokedex) - Lists all caught pokemon in your pokedex