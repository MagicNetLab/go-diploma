package order

import (
	"errors"
)

var ErrorOrderAlreadyAddedByUser = errors.New("order already added by user")
var ErrorOrderAlreadyAddedByOtherUser = errors.New("order already added by other user")

type Order struct {
	Number     int
	Status     string
	Accrual    float32
	UploadedAt string
}

type UserOrdersResponse []Order

type UserBalance struct {
	Current   float64
	Withdrawn float64
}

type Withdraw struct {
	Order       string
	Sum         float64
	ProcessedAt string
}

type WithdrawList []Withdraw
