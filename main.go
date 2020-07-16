package main

import (
	"fmt"
	"log"

	"github.com/lafriks/go-tiled"
)

const (
	mapPath = "arch1.tmx"
)

func main() {
	gameMap, err := tiled.LoadFromFile(mapPath)
	if err != nil {
		log.Fatalf("Failed to load map: %s", err)
	}

	fmt.Printf("%v", gameMap.ObjectGroups[0].Objects[0].X)
}
