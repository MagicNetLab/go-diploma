package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/MagicNetLab/go-diploma/internal/config/env"
	"github.com/MagicNetLab/go-diploma/internal/config/flags"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"go.uber.org/zap"
)

var Env Environment

func GetAppConfig() (AppEnvironment, error) {
	if Env.isValid() {
		return &Env, nil
	}

	getEnvValues()
	getFlagsValues()

	if Env.GetJWTSecret() == "" {
		err := Env.SetJWTSecret(getRandomSecret())
		if err != nil {
			logger.Error("Failed to set JWT secret", zap.String("error", err.Error()))
		}
	}

	if Env.isValid() {
		return &Env, nil
	}

	return &Environment{}, errors.New("invalid config")
}

func getEnvValues() {
	envValues, err := env.Parse()
	if err != nil {
		logger.Error("fail parse env params", zap.Error(err))
		return
	}

	if envValues.HasRunAddress() {
		err = Env.SetRunAddress(envValues.GetRunAddress())
		if err != nil {
			logger.Error("fail set RunAddress from env params", zap.Error(err))
		}
	}

	if envValues.HasDBUri() {
		err = Env.SetDBUri(envValues.GetDBUri())
		if err != nil {
			logger.Error("fail set DBUri from env params", zap.Error(err))
		}
	}

	if envValues.HasAccrualSystemURL() {
		err = Env.SetAccrualSystemURL(envValues.GetAccrualSystemURL())
		if err != nil {
			logger.Error("fail set AccrualSystemUrl from env params", zap.Error(err))
		}
	}

	if envValues.HasJWTSecret() {
		err = Env.SetJWTSecret(envValues.GetJWTSecret())
		if err != nil {
			logger.Error("fail set JWTSecret from env params: %v", zap.Error(err))
		}
	}
}

func getFlagsValues() {
	flagValues, err := flags.Parse()
	if err != nil {
		logger.Error("fail parse flag params", zap.Error(err))
	}

	if flagValues.HasRunAddress() {
		err = Env.SetRunAddress(flagValues.GetRunAddress())
		if err != nil {
			logger.Error("fail set RunAddress from flag params: %v", zap.Error(err))
		}
	}

	if flagValues.HasDBUri() {
		err = Env.SetDBUri(flagValues.GetDBUri())
		if err != nil {
			logger.Error("fail set DBUri from flag params", zap.Error(err))
		}
	}

	if flagValues.HasAccrualSystemURL() {
		err = Env.SetAccrualSystemURL(flagValues.GetAccrualSystemURL())
		if err != nil {
			logger.Error("fail set AccrualSystemUrl from flag params", zap.Error(err))
		}
	}

	if flagValues.HasJWTSecret() {
		err = Env.SetJWTSecret(flagValues.GetJWTSecret())
		if err != nil {
			logger.Error("fail set JWTSecret from flag params", zap.Error(err))
		}
	}
}

func getRandomSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
