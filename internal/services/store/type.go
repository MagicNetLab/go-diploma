package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/golang-migrate/migrate/v4"
	pgsql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

var ErrorWithdrawNotUnique = errors.New("withdraw order number already exists")

type Store struct {
	connectString string
}

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

type Order struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	Number     int       `db:"number"`
	Status     string    `db:"status"`
	Accrual    float32   `db:"accrual"`
	UploadedAt time.Time `db:"uploaded_at"`
}

type Withdraw struct {
	Id          int       `db:"id"`
	UserID      int       `db:"user_id"`
	OrderNum    int       `db:"order_num"`
	Sum         float64   `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
}

type WithdrawList []Withdraw

func (s *Store) SetConnectString(str string) error {
	if str == "" {
		return errors.New("store connect string is empty")
	}

	s.connectString = str
	return nil
}

func (s *Store) Migrate() error {
	db, err := sql.Open("postgres", s.connectString)
	if err != nil {
		return errors.New("could not connect to postgres")
	}

	driver, err := pgsql.WithInstance(db, &pgsql.Config{})
	if err != nil {
		return errors.New("could not connect to postgres")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		logger.Error(fmt.Sprintf("fail init db migrate: %v", err))
		return err
	}

	m.Up()

	return nil
}

func (s *Store) Ping(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.connectString)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	err = conn.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
