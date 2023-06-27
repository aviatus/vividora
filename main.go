package main

import (
	api "aviatus/vividora/api"
	store "aviatus/vividora/internal/store"
)

func main() {
	err := store.StartStore()
	if err != nil {
		panic(err)
	}

	err = api.StartServer()
	if err != nil {
		panic(err)
	}
}
