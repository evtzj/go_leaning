package database

import (
	"home_Tutoring_System/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	return err
}

func Migrate() error {
	return DB.AutoMigrate(
		&model.TeacherProfile{},
		&model.User{},
		&model.ChatMessage{},
		&model.Match{},
		&model.TeacherFavorite{},
		&model.Order{},
	)
}
