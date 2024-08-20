package server

import (
	"net/http"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"go.uber.org/zap"
)

func Run(env config.AppEnvironment) {
	router := getRoute()

	err := http.ListenAndServe(env.GetRunAddress(), router)
	if err != nil {
		logger.Fatal("fail run server", zap.Error(err))
		return
	}
}
