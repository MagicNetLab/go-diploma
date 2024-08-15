package server

import (
	"fmt"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"net/http"

	"github.com/MagicNetLab/go-diploma/internal/config"
)

func Run(env config.AppEnvironment) {
	router := getRoute()

	err := http.ListenAndServe(env.GetRunAddress(), router)
	if err != nil {
		logger.Fatal(fmt.Sprintf("fail run server: %v", err), make(map[string]interface{}))
		return
	}
}
