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
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Redirect  string    `json:"redirect" gorm:"not null"`
	Shortened string    `json:"goly" gorm:"unique;not null"`
	Clicked   uint64    `json:"clicked"`
	ExpiredAt time.Time `json:"expired_at" gorm:"not null"`
	UserID    uint64    `json:"user_id" gorm:"not null"`
}
type User struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	Username          string     `json:"username"`
	HashedPassword    string     `json:"hashed_password" gorm:"not null"`
	PasswordChangedAt time.Time  `json:"password_changed_at" gorm:"not null;default:'0001-01-01 00:00:00Z'"`
	CreatedAt         time.Time  `json:"created_at" gorm:"not null;default:now()"`
	ShortUrls         []ShortUrl `json:"short_urls" gorm:"foreignKey:UserID"`
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
}
