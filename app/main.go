package main

import (
	"log"
	"net/http"
	handler "test_task2/domain_methods/handlres"
	"test_task2/infrastructure/config"
	db "test_task2/infrastructure/database"

	_ "test_task2/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	_ "gorm.io/gorm"
)

// @title Songs API
// @version 1.0
// @description API для управления библиотекой песен.
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9000
// @BasePath /
func main() {
	cfg := config.LoadConfig()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := db.Migrate(cfg.DBURL)
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}

	h := &handler.Handler{
		DB:        db,
		Logger:    logger,
		APIServer: cfg.APIServer,
	}

	r := chi.NewRouter()
	r.Get("/songs", h.GetLibrary)            // Получение библиотеки
	r.Get("/songs/{id}/text", h.GetSongText) // Получение текста песни
	r.Delete("/songs/{id}", h.DeleteSong)    // Удаление песни
	r.Put("/songs/{id}", h.UpdateSong)       // Изменение песни
	r.Post("/songs", h.AddSong)              // Добавление новой песни

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	logger.Info("Starting server on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
