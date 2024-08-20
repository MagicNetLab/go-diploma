package store

import (
	"context"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"go.uber.org/zap"
)

var store Store

func Init(env config.AppEnvironment) error {
	err := store.SetConnectString(env.GetDBUri())
	if err != nil {
		logger.Error("fail set db connect param", zap.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = store.Ping(ctx)
	if err != nil {
		logger.Error("fail ping db connect param", zap.Error(err))
		return err
	}

	err = store.Migrate()
	if err != nil {
		logger.Error("fail migrate", zap.Error(err))
		return err
	}

	return nil
}
