package logger

import (
	"net/http"
	"time"
)

func Middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &ResponseData{
			Status: 0,
			Size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		args := map[string]interface{}{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"status":   responseData.Status,
			"duration": duration,
			"size":     responseData.Size,
		}
		log.Info("request info", args)
	}
}
