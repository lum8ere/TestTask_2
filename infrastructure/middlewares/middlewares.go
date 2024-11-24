package middlewares

import (
	"net/http"
	"test_task2/infrastructure/smartContext"
)

func RecoveryMiddleware(appCtx *smartContext.SmartContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					appCtx.Errorf("Recovered from panic: %v", err)

					// Возвращаем ошибку в ответе
					w.WriteHeader(http.StatusInternalServerError)
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(`{"error": "internal server error"}`))
				}
			}()

			// Передаем управление следующему хендлеру
			next.ServeHTTP(w, r)
		})
	}
}
