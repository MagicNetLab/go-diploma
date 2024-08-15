package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	orderStatusNew        = "NEW"
	orderStatusProcessing = "PROCESSING"
	orderStatusInvalid    = "INVALID"
	orderStatusProcessed  = "PROCESSED"

	getOrderByNumberSQL      = "SELECT id,number,user_id,status,accrual,uploaded_at FROM orders WHERE number = $1"
	insertOrderSQL           = "INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)"
	getOrdersByUserIdSQL     = "SELECT id,user_id,number,status,accrual,uploaded_at FROM orders WHERE user_id = $1"
	accrualAmountByUserIdSQL = "SELECT sum(accrual) FROM orders WHERE user_id = $1 and status = 'PROCESSED'"
)

func GetOrderByNumber(number int) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	var order Order
	row := conn.QueryRow(ctx, getOrderByNumberSQL, number)
	err = row.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func CreateOrder(number int, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	res, err := conn.Exec(ctx, insertOrderSQL, number, userID, orderStatusNew)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("failed to insert order")
	}

	return nil
}

func GetOrdersByUserID(userID int) ([]Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	var orders []Order
	rows, err := conn.Query(ctx, getOrdersByUserIdSQL, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return orders, nil
		}

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err = rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			rows.Close()
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func GetAccrualAmountByUserID(userID int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)

	var amount sql.NullFloat64
	err = conn.QueryRow(ctx, accrualAmountByUserIdSQL, userID).Scan(&amount)
	if err != nil {
		return 0, err
	}

	if !amount.Valid {
		return 0, nil
	}

	return amount.Float64, nil
}
