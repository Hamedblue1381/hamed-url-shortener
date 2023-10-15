package model

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
}
type ShortUrl struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	Redirect     string    `json:"redirect" gorm:"not null"`
	Shortened    string    `json:"goly" gorm:"unique;not null;index"`
	Clicked      uint64    `json:"clicked"`
	UserID       uint64    `json:"user_id"`
	User         User      `gorm:"foreignkey:UserID"`
	LastAccessed time.Time `gorm:"not null;"`
}
type User struct {
	ID                uint64    `json:"id" gorm:"primaryKey"`
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password" gorm:"not null"`
	PasswordChangedAt time.Time `json:"password_changed_at" gorm:"not null;default:'0001-01-01 00:00:00Z'"`
	CreatedAt         time.Time `json:"created_at" gorm:"not null;default:now()"`
}

func Setup(config Config) {
	dsn := config.DBSource
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&ShortUrl{}, &User{})
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			cleanupStaleURLs(db)
			time.Sleep(24 * time.Hour) // Run the cleanup task once a day
		}
	}()

}
func cleanupStaleURLs(db *gorm.DB) {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	db.Where("last_accessed < ?", oneYearAgo).Delete(&ShortUrl{})
}
