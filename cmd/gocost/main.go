package main

import (
	"log"

	"github.com/madalinpopa/gocost/internal/config"
)

func main() {

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

}
