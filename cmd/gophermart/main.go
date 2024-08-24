package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/accrual"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/server"
	"github.com/MagicNetLab/go-diploma/internal/services/store"
)

func main() {
	initApp()
	runServer()
	waitStop()
}

func initApp() {
	err := logger.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	cnf, err := config.GetAppConfig()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Fatal("fail loading config", args)
		return
	}

	err = store.Init(cnf)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Fatal("failed initializing store", args)
		return
	}
}

func runServer() {
	cnf, err := config.GetAppConfig()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Fatal("fail loading config", args)
		return
	}

	// run accrual worker
	accrual.RunWorker()

	// run server
	go func() { server.Run(cnf) }()
}

func waitStop() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownCh

	logger.Info("Shutting down...", nil)
	logger.Sync()
}
