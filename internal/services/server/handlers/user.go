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
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		var regRequest RegisterUserRequest
		if err := json.NewDecoder(r.Body).Decode(&regRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if regRequest.Login == "" || regRequest.Password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err := user.Register(regRequest.Login, regRequest.Password)
		if err != nil {
			if errors.As(err, &user.ErrorUserExists) {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}

			logger.Error(fmt.Sprintf("fail register user: %v", err), make(map[string]interface{}))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token, err := user.Login(regRequest.Login, regRequest.Password)
		if err != nil {
			logger.Error(fmt.Sprintf("fail login: %v", err), make(map[string]interface{}))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Minute * 60)}
		r.AddCookie(&newCookie)
		http.SetCookie(w, &newCookie)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to write response %v", err), make(map[string]interface{}))
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if loginRequest.Login == "" || loginRequest.Password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		token, err := user.Login(loginRequest.Login, loginRequest.Password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		newCookie := http.Cookie{Name: "token", Value: token, Path: "/", Expires: time.Now().Add(time.Minute * 60)}
		r.AddCookie(&newCookie)
		http.SetCookie(w, &newCookie)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("OK"))
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to write response %v", err), make(map[string]interface{}))
		}

	}
}
