package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	hasUserByLoginSQL = "SELECT count(id) FROM users where login = $1"
	insertUserSQL     = "INSERT INTO users (login, password) VALUES ($1, $2)"
	getUserByLoginSQL = "SELECT id,login,password from users where login = $1"
)

func HasUserByLogin(login string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to database: %v", err))
		return false, err
	}
	defer conn.Close(ctx)

	var count int
	err = conn.QueryRow(ctx, hasUserByLoginSQL, login).Scan(&count)
	if err != nil {
		logger.Error(fmt.Sprintf("failed execute query: %v", err))
		return false, err
	}

	return count > 0, nil
}

func CreateUser(login string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to database: %v", err))
		return err
	}
	defer conn.Close(ctx)

	res, err := conn.Exec(ctx, insertUserSQL, login, password)
	if err != nil {
		logger.Error(fmt.Sprintf("failed execute query: %v", err))
		return err
	}

	if res.RowsAffected() == 0 {
		logger.Error("failed execute create user query")
		return errors.New("failed to insert user")
	}

	return nil
}

func GetUserByLogin(login string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, store.connectString)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to connect to database: %v", err))
		return User{}, err
	}
	defer conn.Close(ctx)

	var user User
	err = conn.QueryRow(ctx, getUserByLoginSQL, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		logger.Error(fmt.Sprintf("failed execute query: %v", err))
		return User{}, err
	}

	return user, nil
}
