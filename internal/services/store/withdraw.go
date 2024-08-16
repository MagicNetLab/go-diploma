package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"strings"
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
		logger.Error(fmt.Sprintf("error connecting to database: %v", err))
		return 0, err
	}
	defer conn.Close(ctx)

	var amount sql.NullFloat64
	err = conn.QueryRow(ctx, withdrawAmountByUserIdSQL, userID).Scan(&amount)
	if err != nil {
		logger.Error(fmt.Sprintf("error execute query getting withdraw amount: %v", err))
		return 0, err
	}

	if !amount.Valid {
		logger.Error(fmt.Sprintf("error execute query getting withdraw amount: %v", err))
		return 0, nil
	}

	return amount.Float64, nil
}

func CreateWithdraw(order int, amount float64, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		logger.Error(fmt.Sprintf("error connecting to database: %v", err))
		return err
	}
	defer conn.Close(ctx)

	result, err := conn.Exec(ctx, createWithdrawSQL, order, amount, userID)
	if err != nil {
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			logger.Error(fmt.Sprintf("error creating withdraw: number %v already exists", order))
			return ErrorWithdrawNotUnique
		}

		logger.Error(fmt.Sprintf("failed insert new withdraw: %v", err))
		return err
	}

	if result.RowsAffected() == 0 {
		logger.Error(fmt.Sprintf("error creating withdraw: %v", err))
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
		logger.Error(fmt.Sprintf("error connecting to database: %v", err))
		return list, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, withdrawListByUserIdSQL, userID)
	if err != nil {
		logger.Error(fmt.Sprintf("error execute query getting withdraw list: %v", err))
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

	return list, nil
}
