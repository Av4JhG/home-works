package internalhttp

import (
	"net/http"
	"time"

	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(next http.Handler, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		path := r.URL.Path
		httpVersion := r.Proto
		clientIP := r.RemoteAddr
		currentTime := time.Now().Format(time.RFC3339)
		userAgent := r.UserAgent()
		logger.Infof("Received request: %s %s %s from %s at %s, User-Agent: %s",
			method,
			path,
			httpVersion,
			clientIP,
			currentTime,
			userAgent,
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		logger.Infof("Processed request: %s %s, latency: %s", method, path, duration)
	})
}
