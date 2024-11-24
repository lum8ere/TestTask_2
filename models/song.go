package models

import "time"

// Song структура данных для песни
type Song struct {
	ID          uint       `gorm:"primaryKey" json:"id"`               // Первичный ключ
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`   // Время создания
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`   // Время обновления
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`  // Soft-delete
	Group       string     `gorm:"column:group;not null" json:"group"` // Название группы
	Title       string     `gorm:"column:title;not null" json:"title"` // Название песни
	ReleaseDate string     `gorm:"column:release_date;not null"`       // Дата релиза
	Text        string     `gorm:"column:text" json:"text"`            // Текст песни
	Link        string     `gorm:"column:link" json:"link"`            // Ссылка
}
