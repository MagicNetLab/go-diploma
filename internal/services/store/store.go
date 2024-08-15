package store

import (
	"context"
	"fmt"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/config"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
)

var store Store

func Init(env config.AppEnvironment) error {
	err := store.SetConnectString(env.GetDBUri())
	if err != nil {
		logger.Error(fmt.Sprintf("fail set db conect param: %s", err), make(map[string]interface{}))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = store.Ping(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("fail ping db conect param: %v", err), make(map[string]interface{}))
		return err
	}

	err = store.Migrate()
	if err != nil {
		logger.Error(fmt.Sprintf("fail migrate db conect param: %s", err), make(map[string]interface{}))
		return err
	}

	return nil
}
