package loaders

import (
	"encoding/json"
	"log"
	"os"

	"apodeiktikos.com/fbtest/core"
)

func LoadJson(path string, sprites *core.Sprites) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error opening file %s\n%v", path, err)
	}

	err = json.Unmarshal(data, &sprites)
	if err != nil {
		log.Fatalf("Error loading sprite definition: %s\n%v", path, err)
	}
}
