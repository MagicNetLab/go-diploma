package server

import (
	"net/http"

	"github.com/MagicNetLab/go-diploma/internal/services/compression"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/user"
)

type mwList []func(handlerFunc http.HandlerFunc) http.HandlerFunc

var mwGuestGet = mwList{mwDefault, mwGet, mwGuest}
var mwGuestPost = mwList{mwDefault, mwPost, mwGuest}
var mwAuthorizedGet = mwList{mwDefault, mwGet, mwAuthorized}
var mwAuthorizedPost = mwList{mwDefault, mwPost, mwAuthorized}

func mwDefault(h http.HandlerFunc) http.HandlerFunc {
	return compression.GzipMiddleware(logger.Middleware(h))
}

func mwGuest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user.CheckAuthorize(r) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func mwAuthorized(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !user.CheckAuthorize(r) {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		h.ServeHTTP(w, r)
	}
}

func mwGet(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func mwPost(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func mw(h http.HandlerFunc, mw mwList) http.HandlerFunc {
	f := h
	for _, m := range mw {
		f = m(f)
	}

	return f
}
