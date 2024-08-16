package main

import (
	"github.com/MagicNetLab/go-diploma/internal/services/store"
	"go.uber.org/zap"
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
		logger.Fatal("fail loading config", zap.String("error", err.Error()))
		return
	}

	err = store.Init(cnf)
	if err != nil {
		logger.Fatal("failed initializing store", zap.String("error", err.Error()))
		return
	}

	// run server
	go func() {
		server.Run(cnf)
	}()

	// run workers

	select {}
}
