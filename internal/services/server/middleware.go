package server

import (
	"net/http"

	"github.com/MagicNetLab/go-diploma/internal/services/compression"
	"github.com/MagicNetLab/go-diploma/internal/services/logger"
	"github.com/MagicNetLab/go-diploma/internal/services/user"
)

func mwDefault(h http.HandlerFunc) http.HandlerFunc {
	return compression.GzipMiddleware(logger.Middleware(h))
}

func mwGuest(h http.HandlerFunc) http.HandlerFunc {
	return mwDefault(func(w http.ResponseWriter, r *http.Request) {
		if user.CheckAuthorize(r) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func mwAuthorized(h http.HandlerFunc) http.HandlerFunc {
	return mwDefault(func(w http.ResponseWriter, r *http.Request) {
		if !user.CheckAuthorize(r) {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		h.ServeHTTP(w, r)
	})
}
