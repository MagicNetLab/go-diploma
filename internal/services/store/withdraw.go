package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
)

const (
	withdrawAmountByUserIDSQL = "SELECT sum(sum) FROM withdraw WHERE user_id = $1"
	createWithdrawSQL         = "INSERT INTO withdraw (order_num, sum, user_id) VALUES ($1, $2, $3)"
	withdrawListByUserIDSQL   = "SELECT order_num, sum, processed_at FROM withdraw WHERE user_id = $1 ORDER BY processed_at"
)

func GetWithdrawAmountByUserID(userID int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error connecting to database", args)
		return 0, err
	}
	defer conn.Close(ctx)

	var amount sql.NullFloat64
	err = conn.QueryRow(ctx, withdrawAmountByUserIDSQL, userID).Scan(&amount)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error execute query getting withdraw amount", args)
		return 0, err
	}

	if !amount.Valid {
		args := map[string]interface{}{"error": "amount not found"}
		logger.Error("error execute query getting withdraw amount", args)
		return 0, nil
	}

	return amount.Float64, nil
}

func CreateWithdraw(order int, amount float64, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error connecting to database", args)
		return err
	}
	defer conn.Close(ctx)

	result, err := conn.Exec(ctx, createWithdrawSQL, order, amount, userID)
	if err != nil {
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			args := map[string]interface{}{"number": order}
			logger.Error("error creating withdraw: number already exists", args)
			return ErrorWithdrawNotUnique
		}

		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed insert new withdraw", args)
		return err
	}

	if result.RowsAffected() == 0 {
		args := map[string]interface{}{"number": order}
		logger.Error("error creating withdraw", args)
		return errors.New("failed to insert withdraw")
	}

	return nil
}

func GetWithdrawListByUserID(userID int) (WithdrawList, error) {
	var list WithdrawList

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error connecting to database", args)
		return list, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, withdrawListByUserIDSQL, userID)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error execute query getting withdraw list", args)
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
