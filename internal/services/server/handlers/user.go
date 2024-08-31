package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/user"
)

func UserRegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var regRequest RegisterUserRequest
		if err := json.NewDecoder(r.Body).Decode(&regRequest); err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error decoding register request body", args)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !regRequest.IsValid() {
			logger.Error("failed register user: one or more request params is empty", nil)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := user.Register(regRequest.Login, regRequest.Password)
		if err != nil {
			if errors.Is(err, user.ErrorUserExists) {
				logger.Error(fmt.Sprintf("failed register user: login %s is exists", regRequest.Login), nil)
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			args := map[string]interface{}{"error": err.Error()}
			logger.Error("fail register user", args)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token, err := user.Login(regRequest.Login, regRequest.Password)
		if err != nil {
			args := map[string]interface{}{"error": err.Error(), "login": regRequest.Login}
			logger.Error("fail user login", args)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Minute * 60)}
		r.AddCookie(&newCookie)
		http.SetCookie(w, &newCookie)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("Failed to write register user response response", args)
		}
	}
}

func UserLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("error decoding login request body", args)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !loginRequest.IsValid() {
			logger.Error("failed login user: one or more request params is empty", nil)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		token, err := user.Login(loginRequest.Login, loginRequest.Password)
		if err != nil {
			args := map[string]interface{}{"error": err.Error(), "login": loginRequest.Login}
			logger.Error("fail user login", args)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Minute * 60)}
		r.AddCookie(&newCookie)
		http.SetCookie(w, &newCookie)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("Failed to write response", args)
		}
	}
}
