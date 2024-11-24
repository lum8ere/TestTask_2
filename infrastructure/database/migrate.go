package db

import (
	"test_task2/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Song{})
	if err != nil {
		return nil, err
	}

	// Добавление тестовых данных
	var count int64
	db.Model(&models.Song{}).Count(&count)
	if count == 0 { // Если таблица пустая
		songs := []models.Song{
			{Group: "Muse", Title: "Supermassive Black Hole", ReleaseDate: "2006-07-16", Text: "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight", Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Linkin Park", Title: "In the End", ReleaseDate: "2000-10-24", Text: "It starts with one thing\nI don't know why\nIt doesn't even matter how hard you try\nKeep that in mind...", Link: "https://www.youtube.com/watch?v=eVTXPUF4Oz4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Queen", Title: "Bohemian Rhapsody", ReleaseDate: "1975-10-31", Text: "Is this the real life? Is this just fantasy?\nCaught in a landslide, no escape from reality...", Link: "https://www.youtube.com/watch?v=fJ9rUzIMcZQ", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Nirvana", Title: "Smells Like Teen Spirit", ReleaseDate: "1991-09-10", Text: "Load up on guns, bring your friends\nIt's fun to lose and to pretend...", Link: "https://www.youtube.com/watch?v=hTWKbfoikeg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Adele", Title: "Hello", ReleaseDate: "2015-10-23", Text: "Hello, it's me\nI was wondering if after all these years you'd like to meet...", Link: "https://www.youtube.com/watch?v=YQHsXMglC9A", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "The Beatles", Title: "Hey Jude", ReleaseDate: "1968-08-26", Text: "Hey Jude, don't make it bad\nTake a sad song and make it better...", Link: "https://www.youtube.com/watch?v=A_MjCqQoLLA", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Coldplay", Title: "Fix You", ReleaseDate: "2005-09-05", Text: "When you try your best but you don't succeed\nWhen you get what you want but not what you need...", Link: "https://www.youtube.com/watch?v=k4V3Mo61fJM", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Eminem", Title: "Lose Yourself", ReleaseDate: "2002-10-22", Text: "Look, if you had one shot, or one opportunity\nTo seize everything you ever wanted...", Link: "https://www.youtube.com/watch?v=_Yhyp-_hX2s", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Michael Jackson", Title: "Billie Jean", ReleaseDate: "1983-01-02", Text: "She was more like a beauty queen\nFrom a movie scene...", Link: "https://www.youtube.com/watch?v=Zi_XLOBDo_Y", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Imagine Dragons", Title: "Believer", ReleaseDate: "2017-02-01", Text: "First things first\nI'ma say all the words inside my head...", Link: "https://www.youtube.com/watch?v=7wtfhZwyrcc", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "The Rolling Stones", Title: "Paint It Black", ReleaseDate: "1966-05-07", Text: "I see a red door and I want it painted black\nNo colors anymore I want them to turn black...", Link: "https://www.youtube.com/watch?v=O4irXQhgMqg", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Pink Floyd", Title: "Comfortably Numb", ReleaseDate: "1979-11-23", Text: "Hello,\nIs there anybody in there?\nJust nod if you can hear me\nIs there anyone at home?", Link: "https://www.youtube.com/watch?v=_FrOQC-zEog", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "U2", Title: "With Or Without You", ReleaseDate: "1987-03-16", Text: "See the stone set in your eyes\nSee the thorn twist in your side\nI'll wait for you...", Link: "https://www.youtube.com/watch?v=XmSdTa9kaiQ", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Taylor Swift", Title: "Shake It Off", ReleaseDate: "2014-08-18", Text: "I stay out too late\nGot nothing in my brain\nThat's what people say...", Link: "https://www.youtube.com/watch?v=nfWlot6h_JM", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Group: "Ed Sheeran", Title: "Shape of You", ReleaseDate: "2017-01-06", Text: "The club isn't the best place to find a lover\nSo the bar is where I go...", Link: "https://www.youtube.com/watch?v=JGwWNGJdvx8", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		db.Create(&songs)
	}

	return db, nil
}
