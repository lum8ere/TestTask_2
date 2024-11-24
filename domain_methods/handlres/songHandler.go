package handlers

import (
	"errors"
	"net/http"
	"strings"
	"test_task2/domain_methods/service"
	"test_task2/domain_methods/utils"
	"test_task2/infrastructure/smartContext"
	"test_task2/models"
)

// @Summary Add a new song
// @Description Добавление новой песни в библиотеку
// @Tags Songs
// @Param data body models.Song true "Данные новой песни"
// @Success 201 {object} models.Song
// @Failure 400 {object} map[string]string
// @Router /songs [post]
func AddSongHandler(ctx *smartContext.SmartContext, w http.ResponseWriter, r *http.Request, params map[string]interface{}) (interface{}, error) {
	group, ok := params["group"].(string)
	if !ok {
		ctx.Errorf("Missing 'group' parameter")
		return nil, errors.New("missing group")
	}
	song, ok := params["song"].(string)
	if !ok {
		ctx.Errorf("Missing 'song' parameter")
		return nil, errors.New("missing song")
	}

	ctx.Debugf("Fetching song details for group='%s', song='%s'", group, song)
	details, err := service.FetchSongDetails(ctx.GetAPIServer(), group, song)
	if err != nil {
		return nil, err
	}

	newSong := models.Song{
		Group:       group,
		Title:       song,
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	if err := ctx.GetDB().Create(&newSong).Error; err != nil {
		ctx.Errorf("Failed to save new song: %v", err)
		return nil, err
	}

	ctx.Infof("Successfully added new song: %v", newSong)
	return newSong, nil
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
func GetLibraryHandler(ctx *smartContext.SmartContext, w http.ResponseWriter, r *http.Request, params map[string]interface{}) (interface{}, error) {
	query := ctx.GetDB().Model(&models.Song{})

	// Фильтрация
	if group, ok := params["group"].(string); ok {
		ctx.Errorf("Missing 'group' parameter")
		query = query.Where(`"group" = ?`, group)
	}
	if title, ok := params["title"].(string); ok {
		ctx.Errorf("Missing 'title' parameter")
		query = query.Where("title = ?", title)
	}

	// Пагинация
	page := utils.ParseInt(params["page"], 1)
	limit := utils.ParseInt(params["limit"], 10)
	offset := (page - 1) * limit

	ctx.Debugf("pagination params: page: %v, song: %v, offset: %v", page, limit)

	var songs []models.Song
	if err := query.Offset(offset).Limit(limit).Find(&songs).Error; err != nil {
		return nil, err
	}

	ctx.Infof("Successfully received songs in quantity: %v", len(songs))
	return songs, nil
}

// @Summary Get song text
// @Description Получение текста песни с пагинацией по куплетам
// @Tags Songs
// @Param id path int true "ID песни"
// @Param page query int true "Номер куплета" default(1)
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id}/text [get]
func GetSongTextHandler(ctx *smartContext.SmartContext, w http.ResponseWriter, r *http.Request, params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		ctx.Errorf("Missing 'id' parameter")
		return nil, errors.New("missing id")
	}
	page := utils.ParseInt(params["page"], 1)

	var song models.Song
	if err := ctx.GetDB().First(&song, id).Error; err != nil {
		return nil, errors.New("song not found")
	}
	ctx.Debugf("Got a song: %v", song.Title)

	verses := strings.Split(song.Text, "\n\n")
	if page < 1 || page > len(verses) {
		return nil, errors.New("invalid page number")
	}

	ctx.Infof("verses on the way out: %v", verses[page-1])
	return map[string]string{"verse": verses[page-1]}, nil
}

// @Summary Delete a song
// @Description Удаление песни из библиотеки
// @Tags Songs
// @Param id path int true "ID песни"
// @Success 204 "Песня успешно удалена"
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [delete]
func DeleteSongHandler(ctx *smartContext.SmartContext, w http.ResponseWriter, r *http.Request, params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		ctx.Errorf("Missing 'id' parameter")
		return nil, errors.New("missing id")
	}

	if err := ctx.GetDB().Delete(&models.Song{}, id).Error; err != nil {
		ctx.Errorf("Failed to delete song: %v", err)
		return nil, err
	}

	ctx.Infof("Deleted song with ID: %v", id)
	return map[string]string{"message": "song deleted successfully"}, nil
}

// @Summary Update a song
// @Description Обновление данных песни
// @Tags Songs
// @Param id path int true "ID песни"
// @Param data body models.Song true "Данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {object} map[string]string
// @Router /songs/{id} [put]
func UpdateSongHandler(ctx *smartContext.SmartContext, w http.ResponseWriter, r *http.Request, params map[string]interface{}) (interface{}, error) {
	id, ok := params["id"].(string)
	if !ok {
		ctx.Errorf("Missing 'id' parameter")
		return nil, errors.New("missing id")
	}

	var song models.Song
	if err := ctx.GetDB().First(&song, id).Error; err != nil {
		ctx.Errorf("Failed to find song: %v", err)
		return nil, errors.New("song not found")
	}

	ctx.Debugf("Song: %v", song)

	if group, ok := params["group"].(string); ok {
		song.Group = group
	}
	if title, ok := params["title"].(string); ok {
		song.Title = title
	}
	if releaseDate, ok := params["release_date"].(string); ok {
		song.ReleaseDate = releaseDate
	}
	if text, ok := params["text"].(string); ok {
		song.Text = text
	}
	if link, ok := params["link"].(string); ok {
		song.Link = link
	}

	if err := ctx.GetDB().Save(&song).Error; err != nil {
		return nil, err
	}

	ctx.Infof("Successfully updated song: %v", song)
	return song, nil
}
