package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"synthori/ediary/m/handlers"
	"time"
)

//	logs {
//			method
//			request
//			body
//			remoteaddr
//			deviceinfo
//	}
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)

		// Собираем fullUrl с корректной проверкой обеих переменных окружения
		fullUrl := ""
		requestUrl := r.URL.RequestURI()
		defaultUrl, urlExists := os.LookupEnv("DEFAULT_URL")
		defaultPort, portExists := os.LookupEnv("DEFAULT_PORT")
		if urlExists && portExists {
			fullUrl = defaultUrl + ":" + defaultPort + requestUrl
		}

		// Получаем данные об устройстве
		device := handlers.GetDeviceInfo(r)

		// Логируем все в ОДНУ JSON-строку
		slog.Info("http request processed",
			"method", r.Method,
			"request_url", fullUrl,
			"user_agent", device.UserAgent,
			"remote_addr", device.IPAddress,
			"duration_ms", time.Since(start).Milliseconds(), // Удобнее для графиков и аналитики в миллисекундах
		)
	})
}
