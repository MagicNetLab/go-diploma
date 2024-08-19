package accrual

import "github.com/go-resty/resty/v2"

const (
	queueSize          = 10
	accrualServicePath = "/api/orders/%s"
)

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

var orderCh chan string
var pauseCh chan string
var serviceHost string
var doneCh chan struct{}

var httpc *resty.Client
