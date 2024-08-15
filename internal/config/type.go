package config

import "errors"

type AppEnvironment interface {
	isValid() bool
	SetRunAddress(addr string) error
	GetRunAddress() string
	SetDBUri(uri string) error
	GetDBUri() string
	SetAccrualSystemUrl(url string) error
	GetAccrualSystemUrl() string
	SetJWTSecret(secret string) error
	GetJWTSecret() string
}

type Environment struct {
	runAddress       string
	dbUri            string
	accrualSystemUri string
	jwtSecret        string
}

func (e *Environment) isValid() bool {
	return e.runAddress != "" && e.dbUri != "" && e.accrualSystemUri != ""
}

func (e *Environment) SetRunAddress(addr string) error {
	if addr == "" {
		return errors.New("fail set RunAddress: value is empty")
	}

	e.runAddress = addr
	return nil
}

func (e *Environment) GetRunAddress() string {
	return e.runAddress
}

func (e *Environment) SetDBUri(uri string) error {
	if uri == "" {
		return errors.New("fail set DBUri: value is empty")
	}

	e.dbUri = uri
	return nil
}

func (e *Environment) GetDBUri() string {
	return e.dbUri
}

func (e *Environment) SetAccrualSystemUrl(url string) error {
	if url == "" {
		return errors.New("fail set AccrualSystemUri: value is empty")
	}
	e.accrualSystemUri = url
	return nil
}

func (e *Environment) GetAccrualSystemUrl() string {
	return e.accrualSystemUri
}

func (e *Environment) SetJWTSecret(secret string) error {
	if secret == "" {
		return errors.New("fail set JWTSecret: value is empty")
	}

	e.jwtSecret = secret
	return nil
}

func (e *Environment) GetJWTSecret() string {
	return e.jwtSecret
}
