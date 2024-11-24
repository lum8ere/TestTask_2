package main

import (
	"net/http"
	handlers "test_task2/domain_methods/handlres"
	"test_task2/domain_methods/utils"
	"test_task2/infrastructure/config"
	db "test_task2/infrastructure/database"
	"test_task2/infrastructure/logger"
	"test_task2/infrastructure/smartContext"

	_ "test_task2/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

	// Создаем логгер
	log := logger.NewLogger()

	// Создаем подключение к базе данных
	database, err := db.Migrate(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Создаем AppContext
	ctx := smartContext.NewSmartContext(database, log, cfg.APIServer)

	// Настройка маршрутов
	r := chi.NewRouter()
	r.Get("/songs", utils.HandleWrapper(ctx, handlers.GetLibraryHandler, "group", "title", "page", "limit"))
	r.Post("/songs", utils.HandleWrapper(ctx, handlers.AddSongHandler, "group", "song"))
	r.Get("/songs/{id}/text", utils.HandleWrapper(ctx, handlers.GetSongTextHandler, "id", "page"))
	r.Delete("/songs/{id}", utils.HandleWrapper(ctx, handlers.DeleteSongHandler, "id"))
	r.Put("/songs/{id}", utils.HandleWrapper(ctx, handlers.UpdateSongHandler, "id", "group", "title", "release_date", "text", "link"))

	// Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Infof("Starting server on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
