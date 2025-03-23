package main

import (
	"log"

	"github.com/dnagikh/gockuper-cli/cmd"
	"github.com/dnagikh/gockuper-cli/config"
	"github.com/dnagikh/gockuper-cli/internal/logger"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = logger.InitLogger()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}

	cmd.Execute()
}
