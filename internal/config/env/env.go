package env

import (
	"fmt"
	"os"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/joho/godotenv"
)

func Parse() (Options, error) {
	var opts Options

	err := godotenv.Load(".env")
	if err != nil {
		logger.Error(fmt.Sprintf("fail load .env file: %s", err))
	}

	runAddressValue := os.Getenv(runAddressKey)
	if runAddressValue != "" {
		opts.runAddress = runAddressValue
	}

	dbIriValue := os.Getenv(dbUriKey)
	if dbIriValue != "" {
		opts.dbUri = dbIriValue
	}

	accrualSystemUrlValue := os.Getenv(accrualSystemUrlKey)
	if accrualSystemUrlValue != "" {
		opts.accrualSystemUrl = accrualSystemUrlValue
	}

	jwtSecret := os.Getenv(jwtSecret)
	if jwtSecret != "" {
		opts.jwtSecret = jwtSecret
	}

	return opts, nil
}
