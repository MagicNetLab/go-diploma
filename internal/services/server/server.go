package server

import (
	"fmt"
	"net/http"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
)

func Run(env config.AppEnvironment) {
	router := getRoute()

	err := http.ListenAndServe(env.GetRunAddress(), router)
	if err != nil {
		logger.Fatal(fmt.Sprintf("fail run server: %v", err))
		return
	}
}
