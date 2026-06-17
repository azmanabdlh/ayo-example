package main

import (
	"log"

	"github.com/azmanabdlh/ayo-example/internal/config"
)

func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to init the config: %v", err)
	}

	err = startApp(cfg)
	if err != nil {
		log.Fatalf("failed to start the app: %v", err)
	}
}
