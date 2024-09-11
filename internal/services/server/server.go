package server

import (
	"net/http"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
)

func Run(env config.AppEnvironment) {
	router := getRoute()

	err := http.ListenAndServe(env.GetRunAddress(), router)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Fatal("fail run server", args)
		return
	}
}
