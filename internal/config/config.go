package config

import (
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

	if Env.isValid() {
		return &Env, nil
	}

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
}
