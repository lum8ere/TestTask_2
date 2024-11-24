package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"test_task2/domain_methods/service"
	"test_task2/domain_methods/utils"
	"test_task2/models"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	Logger    *zap.Logger
	APIServer string
}

// @Summary Add a new song
// @Description Добавление новой песни в библиотеку
// @Tags Songs
// @Param data body models.Song true "Данные новой песни"
// @Success 201 {object} models.Song
// @Failure 400 {object} map[string]string
// @Router /songs [post]
func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Logger.Error("Invalid input", zap.Error(err))
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	details, err := service.FetchSongDetails(h.APIServer, input.Group, input.Song)
	if err != nil {
		h.Logger.Error("Failed to fetch song details", zap.Error(err))
		http.Error(w, "Failed to fetch song details", http.StatusInternalServerError)
		return
	}

	song := models.Song{
		Group:       input.Group,
		Title:       input.Song,
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	if err := h.DB.Create(&song).Error; err != nil {
		h.Logger.Error("Failed to save song", zap.Error(err))
		http.Error(w, "Failed to save song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// @Summary Get library of songs
// @Description Получение библиотеки песен с фильтрацией по полям и пагинацией
// @Tags Songs
// @Param group query string false "Название группы"
// @Param title query string false "Название песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Success 200 {array} models.Song
// @Failure 400 {object} map[string]string
// @Router /songs [get]
func (h *Handler) GetLibrary(w http.ResponseWriter, r *http.Request) {
	query := h.DB.Model(&models.Song{})
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// Фильтрация
	if group := r.URL.Query().Get("group"); group != "" {
		query = query.Where(`"group" = ?`, group) // Экранируем "group"
	}
	if title := r.URL.Query().Get("title"); title != "" {
		query = query.Where("title = ?", title)
	}

	// Пагинация
	var songs []models.Song
	p, l := 1, 10 // значения по умолчанию
	if page != "" {
		p = utils.ParseInt(page, 1)
	}
	if limit != "" {
		l = utils.ParseInt(limit, 10)
	}
	offset := (p - 1) * l
	query.Offset(offset).Limit(l).Find(&songs)

	// Ответ
	json.NewEncoder(w).Encode(songs)
}

// @Summary Get song text
// @Description Получение текста песни с пагинацией по куплетам
// @Tags Songs
// @Param id path int true "ID песни"
// @Param page query int true "Номер куплета" default(1)
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id}/text [get]
func (h *Handler) GetSongText(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")
	page := utils.ParseInt(r.URL.Query().Get("page"), 1)

	var song models.Song
	if err := h.DB.First(&song, songID).Error; err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	verses := strings.Split(song.Text, "\n\n")
	if page < 1 || page > len(verses) {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"verse": verses[page-1],
	}
	json.NewEncoder(w).Encode(response)
}

// @Summary Delete a song
// @Description Удаление песни из библиотеки
// @Tags Songs
// @Param id path int true "ID песни"
// @Success 204 "Песня успешно удалена"
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")

	if err := h.DB.Delete(&models.Song{}, songID).Error; err != nil {
		http.Error(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update a song
// @Description Обновление данных песни
// @Tags Songs
// @Param id path int true "ID песни"
// @Param data body models.Song true "Данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {object} map[string]string
// @Router /songs/{id} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "id")

	var input struct {
		Group       *string `json:"group"`
		Title       *string `json:"title"`
		ReleaseDate *string `json:"release_date"`
		Text        *string `json:"text"`
		Link        *string `json:"link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var song models.Song
	if err := h.DB.First(&song, songID).Error; err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	if input.Group != nil {
		song.Group = *input.Group
	}
	if input.Title != nil {
		song.Title = *input.Title
	}
	if input.ReleaseDate != nil {
		song.ReleaseDate = *input.ReleaseDate
	}
	if input.Text != nil {
		song.Text = *input.Text
	}
	if input.Link != nil {
		song.Link = *input.Link
	}

	if err := h.DB.Save(&song).Error; err != nil {
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(song)
}
