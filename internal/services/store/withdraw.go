package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/jackc/pgx/v5"
)

const (
	withdrawAmountByUserIdSQL = "SELECT sum(sum) FROM withdraw WHERE user_id = $1"
	createWithdrawSQL         = "INSERT INTO withdraw (order_num, sum, user_id) VALUES ($1, $2, $3)"
	withdrawListByUserIdSQL   = "SELECT order_num, sum, processed_at FROM withdraw WHERE user_id = $1 ORDER BY processed_at"
)

func GetWithdrawAmountByUserID(userID int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		return 0, err
	}
	defer conn.Close(ctx)

	var amount sql.NullFloat64
	err = conn.QueryRow(ctx, withdrawAmountByUserIdSQL, userID).Scan(&amount)
	if err != nil {
		return 0, err
	}

	if !amount.Valid {
		return 0, nil
	}

	return amount.Float64, nil
}

func CreateWithdraw(order int, amount float64, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	result, err := conn.Exec(ctx, createWithdrawSQL, order, amount, userID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed insert new withdraw: %v", err))
		// todo unique error handle
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("failed to insert withdraw")
	}

	return nil
}

func GetWithdrawListByUserId(userID int) (WithdrawList, error) {
	var list WithdrawList

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		return list, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, withdrawListByUserIdSQL, userID)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdraw Withdraw
		err = rows.Scan(&withdraw.OrderNum, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			rows.Close()
			return list, err
		}
		list = append(list, withdraw)
	}
	//err = rows.Err()
	//if err != nil {
	//	return list, err
	//}

	return list, nil
}
