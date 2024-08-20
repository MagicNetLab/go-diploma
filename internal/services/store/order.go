package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"

	getOrderByNumberSQL      = "SELECT id,number,user_id,status,accrual,uploaded_at FROM orders WHERE number = $1"
	insertOrderSQL           = "INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)"
	getOrdersByUserIDSQL     = "SELECT id,user_id,number,status,accrual,uploaded_at FROM orders WHERE user_id = $1"
	accrualAmountByUserIDSQL = "SELECT sum(accrual) FROM orders WHERE user_id = $1 and status = $2"
	updateOrderSQL           = "UPDATE orders SET status = $1, accrual = $2 WHERE id = $3"
)

func GetOrderByNumber(number int) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var order Order

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))
		return order, err
	}
	defer conn.Close(ctx)

	row := conn.QueryRow(ctx, getOrderByNumberSQL, number)
	err = row.Scan(&order.ID, &order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		logger.Error("failed to query order", zap.Error(err))
		return order, err
	}

	return order, nil
}

func CreateOrder(number int, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))
		return err
	}
	defer conn.Close(ctx)

	res, err := conn.Exec(ctx, insertOrderSQL, number, userID, OrderStatusNew)
	if err != nil {
		logger.Error("failed to insert order", zap.Error(err))
		return err
	}

	if res.RowsAffected() == 0 {
		logger.Error("failed to insert order")
		return errors.New("failed to insert order")
	}

	return nil
}

func GetOrdersByUserID(userID int) ([]Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))
		return nil, err
	}
	defer conn.Close(ctx)

	var orders []Order
	rows, err := conn.Query(ctx, getOrdersByUserIDSQL, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return orders, nil
		}
		logger.Error("failed to query orders", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err = rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			logger.Error("failed to scan query order", zap.Error(err))
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
		logger.Error("failed to connect to database", zap.Error(err))
		return 0, err
	}
	defer conn.Close(ctx)

	var amount sql.NullFloat64
	err = conn.QueryRow(ctx, accrualAmountByUserIDSQL, userID, OrderStatusProcessed).Scan(&amount)
	if err != nil {
		logger.Error("failed to query accrual amount", zap.Error(err))
		return 0, err
	}

	if !amount.Valid {
		logger.Error("failed to query accrual amount: result is not valid")
		return 0, nil
	}

	return amount.Float64, nil
}

func UpdateOrder(order Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))
		return err
	}
	defer conn.Close(ctx)

	res, err := conn.Exec(ctx, updateOrderSQL, order.Status, order.Accrual, order.ID)
	if err != nil {
		logger.Error("failed to update order", zap.Error(err))
		return err
	}

	if res.RowsAffected() == 0 {
		logger.Error("failed to update order")
		return errors.New("failed to update order")
	}

	return nil
}
