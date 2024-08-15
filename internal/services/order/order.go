package order

import (
	"errors"
	"strconv"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/services/store"
	"github.com/jackc/pgx/v5"
)

func CreateOrder(number int, userID int) error {
	if !checkNumber(number) {
		return errors.New("invalid number")
	}

	order, err := store.GetOrderByNumber(number)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if order != nil {
		if order.UserID != userID {
			return ErrorOrderAlreadyAddedByOtherUser
		}
		return ErrorOrderAlreadyAddedByUser
	}

	err = store.CreateOrder(number, userID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserOrders(userID int) ([]Order, error) {
	orders, err := store.GetOrdersByUserID(userID)
	if err != nil {
		return nil, err
	}

	var result []Order
	for _, val := range orders {
		order := Order{
			Number:     val.Number,
			Status:     val.Status,
			Accrual:    val.Accrual,
			UploadedAt: val.UploadedAt.Format(time.RFC3339),
		}

		result = append(result, order)
	}

	return result, nil
}

func GetBalanceByUserID(userID int) (UserBalance, error) {
	accrualAmount, err := store.GetAccrualAmountByUserID(userID)
	if err != nil {
		return UserBalance{}, err
	}

	withdrawAmount, err := store.GetWithdrawAmountByUserID(userID)
	if err != nil {
		return UserBalance{}, err
	}

	return UserBalance{Current: accrualAmount - withdrawAmount, Withdrawn: withdrawAmount}, nil
}

func CreateWithdraw(number string, amount float64, userID int) error {
	orderNum, err := strconv.Atoi(number)
	if err != nil {
		return err
	}

	if !checkNumber(orderNum) {
		return errors.New("invalid order number")
	}

	err = store.CreateWithdraw(orderNum, amount, userID)
	if err != nil {
		return err
	}

	return nil
}

func GetWithdrawsByUserId(userID int) (WithdrawList, error) {
	var list WithdrawList

	dbResult, err := store.GetWithdrawListByUserId(userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return list, nil
		}

		return list, err
	}

	for _, val := range dbResult {
		w := Withdraw{Order: strconv.Itoa(val.OrderNum), Sum: val.Sum, ProcessedAt: val.ProcessedAt.Format(time.RFC3339)}
		list = append(list, w)
	}

	return list, nil
}

// todo refactoring luhn checker
func checkNumber(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
