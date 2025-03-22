package main

import (
	"github.com/dnagikh/gockuper-cli/cmd"
	"github.com/dnagikh/gockuper-cli/config"
	"github.com/dnagikh/gockuper-cli/internal/logger"
	"log"
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
