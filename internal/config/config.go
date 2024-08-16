package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/MagicNetLab/go-diploma/internal/config/env"
	"github.com/MagicNetLab/go-diploma/internal/config/flags"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
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
			logger.Error(fmt.Sprintf("Failed to set JWT secret: %v", err))
		}
	}

	if Env.isValid() {
		return &Env, nil
	}

	logger.Error(fmt.Sprintf("Failed buld correct app evironment: %v", Env))

	return &Environment{}, errors.New("invalid config")
}

func getEnvValues() {
	envValues, err := env.Parse()
	if err != nil {
		logger.Error(fmt.Sprintf("fail parse env params: %v", err))
		return
	}

	if envValues.HasRunAddress() {
		err = Env.SetRunAddress(envValues.GetRunAddress())
		if err != nil {
			logger.Error(fmt.Sprintf("fail set RunAddress from env params: %v", err))
		}
	}

	if envValues.HasDBUri() {
		err = Env.SetDBUri(envValues.GetDBUri())
		if err != nil {
			logger.Error(fmt.Sprintf("fail set DBUri from env params: %v", err))
		}
	}

	if envValues.HasAccrualSystemUrl() {
		err = Env.SetAccrualSystemUrl(envValues.GetAccrualSystemUrl())
		if err != nil {
			logger.Error(fmt.Sprintf("fail set AccrualSystemUrl from env params: %v", err))
		}
	}

	if envValues.HasJWTSecret() {
		err = Env.SetJWTSecret(envValues.GetJWTSecret())
		if err != nil {
			logger.Error(fmt.Sprintf("fail set JWTSecret from env params: %v", err))
		}
	}
}

func getFlagsValues() {
	flagValues, err := flags.Parse()
	if err != nil {
		logger.Error(fmt.Sprintf("fail parse flag params: %v", err))
	}

	if flagValues.HasRunAddress() {
		err = Env.SetRunAddress(flagValues.GetRunAddress())
		if err != nil {
			logger.Error(fmt.Sprintf("fail set RunAddress from flag params: %v", err))
		}
	}

	if flagValues.HasDBUri() {
		err = Env.SetDBUri(flagValues.GetDBUri())
		if err != nil {
			logger.Error(fmt.Sprintf("fail set DBUri from flag params: %v", err))
		}
	}

	if flagValues.HasAccrualSystemUrl() {
		err = Env.SetAccrualSystemUrl(flagValues.GetAccrualSystemUrl())
		if err != nil {
			logger.Error(fmt.Sprintf("fail set AccrualSystemUrl from flag params: %v", err))
		}
	}

	if flagValues.HasJWTSecret() {
		err = Env.SetJWTSecret(flagValues.GetJWTSecret())
		if err != nil {
			logger.Error(fmt.Sprintf("fail set JWTSecret from flag params: %v", err))
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
