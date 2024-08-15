package main

import (
	"fmt"
	"github.com/MagicNetLab/go-diploma/internal/services/store"
	"log"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/server"
)

func main() {
	err := logger.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	cnf, err := config.GetAppConfig()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error loading config: %v", err), make(map[string]interface{}))
		return
	}

	err = store.Init(cnf)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error initializing store: %v", err), make(map[string]interface{}))
		return
	}

	server.Run(cnf)
}
