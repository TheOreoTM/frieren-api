package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add colored logrus fields to the request
			logger := logger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
			})

			logger.Infof("Request received")

			logger = logger.WithFields(logrus.Fields{
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			})

			defer func() {
				logger.Infof("Request completed")
			}()

			next.ServeHTTP(w, r)
		})
	}
}
