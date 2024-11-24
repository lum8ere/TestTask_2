package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test_task2/infrastructure/smartContext"

	"github.com/go-chi/chi/v5"
)

// HandlerFunc - универсальный тип функции-обработчика
// Возвращает `interface{}`, который будет автоматически сериализован в JSON
type HandlerFunc func(ctx *smartContext.SmartContext, w http.ResponseWriter, r *http.Request, params map[string]interface{}) (interface{}, error)

// HandleWrapper - враппер для хендлеров
func HandleWrapper(ctx *smartContext.SmartContext, handler HandlerFunc, paramKeys ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Собираем параметры из URL и query
		params := make(map[string]interface{})
		for _, key := range paramKeys {
			if value := chi.URLParam(r, key); value != "" {
				params[key] = value
			}
			if value := r.URL.Query().Get(key); value != "" {
				params[key] = value
			}
		}

		// Если в запросе есть тело, декодируем JSON
		if r.Body != nil && (r.Method == http.MethodPost || r.Method == http.MethodPut) {
			var bodyParams map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&bodyParams); err == nil {
				for key, value := range bodyParams {
					params[key] = value
				}
			}
		}

		// Вызываем логику обработчика
		response, err := handler(ctx, w, r, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Если все прошло успешно, возвращаем JSON-ответ
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func ParseInt(value interface{}, defaultValue int) int {
	if str, ok := value.(string); ok {
		if result, err := strconv.Atoi(str); err == nil {
			return result
		}
	}
	if num, ok := value.(int); ok {
		return num
	}
	return defaultValue
}
