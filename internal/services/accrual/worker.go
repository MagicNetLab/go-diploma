package accrual

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/go-resty/resty/v2"
)

func RunWorker() {
	appConf, err := config.GetAppConfig()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error loading application config: %v", err))
		return
	}

	orderCh = make(chan string, queueSize)
	pauseCh = make(chan string)
	serviceHost = appConf.GetAccrualSystemUrl()

	httpc = resty.New().SetBaseURL(fmt.Sprintf("http://%s", serviceHost))

	for i := 0; i < queueSize; i++ {
		go worker()
	}
}

func worker() {
	for {
		select {
		case order := <-orderCh:
			err := checkOrder(order)
			if err != nil {
				logger.Error("fail checking order: " + err.Error())
			}
		case pause := <-pauseCh:
			p, err := strconv.Atoi(pause)
			if err == nil {
				time.Sleep(time.Duration(p) * time.Second)
			} else {
				logger.Error("fail converting pause to int: " + err.Error())
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
