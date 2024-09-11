package user

import (
	"errors"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var ErrorUserExists = errors.New("user already exists")

const tokenLifetime = 60 * time.Minute
