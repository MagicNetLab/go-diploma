package env

import (
	"go.uber.org/zap"
	"os"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/joho/godotenv"
)

func Parse() (Options, error) {
	var opts Options

	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("fail load .env file", zap.String("error", err.Error()))
	}

	runAddressValue := os.Getenv(runAddressKey)
	if runAddressValue != "" {
		opts.runAddress = runAddressValue
	}

	dbIriValue := os.Getenv(dbURIKey)
	if dbIriValue != "" {
		opts.dbURI = dbIriValue
	}

	accrualSystemURLValue := os.Getenv(accrualSystemURLKey)
	if accrualSystemURLValue != "" {
		opts.accrualSystemURL = accrualSystemURLValue
	}

	jwtSecret := os.Getenv(jwtSecret)
	if jwtSecret != "" {
		opts.jwtSecret = jwtSecret
	}

	return opts, nil
}
