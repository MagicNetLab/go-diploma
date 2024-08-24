package store

import (
	"context"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
)

var store Store

func Init(env config.AppEnvironment) error {
	err := store.SetConnectString(env.GetDBUri())
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("fail set db connect param", args)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = store.Ping(ctx)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("fail ping db connect param", args)
		return err
	}

	err = store.Migrate()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("fail migrate", args)
		return err
	}

	return nil
}
