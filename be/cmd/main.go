package main

import (
	"log"

	"github.com/rentiansheng/dashboard/pkg/config"
)

func main() {
	if err := config.NewServerConfig(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := initDependResource(); err != nil {
		log.Fatalf("failed to init depend resource: %v", err)
	}

	if err := initServer(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
