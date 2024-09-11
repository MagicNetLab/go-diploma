package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/MagicNetLab/go-diploma/internal/config/env"
	"github.com/MagicNetLab/go-diploma/internal/config/flags"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
)

var conf Environment

func GetAppConfig() (AppEnvironment, error) {
	if conf.isValid() {
		return &conf, nil
	}

	getEnvValues()
	getFlagsValues()

	if conf.GetJWTSecret() == "" {
		err := conf.SetJWTSecret(getRandomSecret())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("Failed to set JWT secret", args)
		}
	}

	if conf.isValid() {
		return &conf, nil
	}

	return nil, errors.New("invalid config")
}

func getEnvValues() {
	envValues, err := env.Parse()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("fail parse env params", args)
		return
	}

	if envValues.HasRunAddress() {
		err = conf.SetRunAddress(envValues.GetRunAddress())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail set RunAddress from env params", args)
		}
	}

	if envValues.HasDBUri() {
		err = conf.SetDBUri(envValues.GetDBUri())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail set DBUri from env params", args)
		}
	}

	if envValues.HasAccrualSystemURL() {
		err = conf.SetAccrualSystemURL(envValues.GetAccrualSystemURL())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail set AccrualSystemUrl from env params", args)
		}
	}

	if envValues.HasJWTSecret() {
		err = conf.SetJWTSecret(envValues.GetJWTSecret())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail set JWTSecret from env params: %v", args)
		}
	}
}

func getFlagsValues() {
	flagValues, err := flags.Parse()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("fail parse flag params", args)
	}

	if flagValues.HasRunAddress() {
		err = conf.SetRunAddress(flagValues.GetRunAddress())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail set RunAddress from flag params: %v", args)
		}
	}

	if flagValues.HasDBUri() {
		err = conf.SetDBUri(flagValues.GetDBUri())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail set DBUri from flag params", args)
		}
	}

	if flagValues.HasAccrualSystemURL() {
		err = conf.SetAccrualSystemURL(flagValues.GetAccrualSystemURL())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail set AccrualSystemUrl from flag params", args)
		}
	}

	if flagValues.HasJWTSecret() {
		err = conf.SetJWTSecret(flagValues.GetJWTSecret())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail set JWTSecret from flag params", args)
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
