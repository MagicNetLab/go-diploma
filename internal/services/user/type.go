package user

import "errors"

type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var ErrorUserExists = errors.New("user already exists")
