package accrual

import (
	"strconv"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/go-resty/resty/v2"
)

func RunWorker() {
	appConf, err := config.GetAppConfig()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Fatal("Error loading application config", args)
		return
	}

	orderCh = make(chan string, queueSize)
	pauseCh = make(chan string)
	serviceHost = appConf.GetAccrualSystemURL()

	httpc = resty.New().SetBaseURL(serviceHost)

	for i := 0; i < queueSize; i++ {
		go worker()
	}
}

func worker() {
	for {
		select {
		case order := <-orderCh:
			err := processOrderAccrual(order)
			if err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("fail checking order ", args)
			}
		case pause := <-pauseCh:
			p, err := strconv.Atoi(pause)
			if err == nil {
				time.Sleep(time.Duration(p) * time.Second)
			} else {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("fail converting pause to int", args)
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
