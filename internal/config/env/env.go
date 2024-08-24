package env

import (
	"os"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/joho/godotenv"
)

func Parse() (Options, error) {
	var opts Options

	err := godotenv.Load(".env")
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("fail load .env file", args)
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
