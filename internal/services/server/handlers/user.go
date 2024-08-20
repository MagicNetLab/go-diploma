package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/user"
	"go.uber.org/zap"
)

func UserRegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var regRequest RegisterUserRequest
		if err := json.NewDecoder(r.Body).Decode(&regRequest); err != nil {
			logger.Error("error decoding register request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !regRequest.IsValid() {
			logger.Error("failed register user: one or more request params is empty")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := user.Register(regRequest.Login, regRequest.Password)
		if err != nil {
			if errors.Is(err, user.ErrorUserExists) {
				logger.Error(fmt.Sprintf("failed register user: login %s is exists", regRequest.Login))
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			logger.Error("fail register user", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token, err := user.Login(regRequest.Login, regRequest.Password)
		if err != nil {
			logger.Error("fail user login", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Minute * 60)}
		r.AddCookie(&newCookie)
		http.SetCookie(w, &newCookie)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			logger.Error("Failed to write register user response response", zap.Error(err))
		}
	}
}

func UserLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var loginRequest UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			logger.Error("error decoding login request body", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !loginRequest.IsValid() {
			logger.Error("failed login user: one or more request params is empty")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		token, err := user.Login(loginRequest.Login, loginRequest.Password)
		if err != nil {
			logger.Error("fail user login", zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Minute * 60)}
		r.AddCookie(&newCookie)
		http.SetCookie(w, &newCookie)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			logger.Error("Failed to write response", zap.Error(err))
		}

	}
}
