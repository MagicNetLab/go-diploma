package main

import (
	"log"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/accrual"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/server"
	"github.com/MagicNetLab/go-diploma/internal/services/store"
	"go.uber.org/zap"
)

func main() {
	err := logger.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	cnf, err := config.GetAppConfig()
	if err != nil {
		logger.Fatal("fail loading config", zap.Error(err))
		return
	}

	err = store.Init(cnf)
	if err != nil {
		logger.Fatal("failed initializing store", zap.String("error", err.Error()))
		return
	}

	// run accrual worker
	accrual.RunWorker()

	// run server
	go func() {
		server.Run(cnf)
	}()

	select {}
}
