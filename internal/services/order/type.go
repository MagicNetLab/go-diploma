package order

import (
	"errors"
)

var ErrorOrderAlreadyAddedByUser = errors.New("order already added by user")
var ErrorOrderAlreadyAddedByOtherUser = errors.New("order already added by other user")

type Order struct {
	Number     int     `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
}

type UserOrdersResponse []Order

type UserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type WithdrawList []Withdraw
