package accrual

import "github.com/go-resty/resty/v2"

const (
	queueSize          = 3
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

var httpc *resty.Client
