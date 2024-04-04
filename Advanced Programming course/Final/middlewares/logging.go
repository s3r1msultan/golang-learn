package middlewares

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			// Добавьте другие поля по желанию
		}).Info("Received request")

		next.ServeHTTP(w, r) // передаем управление следующему обработчику

		// Здесь можно добавить логирование ответа, если нужно
	})
}
