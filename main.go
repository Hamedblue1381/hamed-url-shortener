package main

import (
	"github.com/HamedBlue1381/hamed-url-shortener/util"
	"github.com/Hamedblue1381/hamed-url-shortener/db/model"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		panic("cannot load config:")
	}
	dbConfig := model.Config{
		DBDriver: config.DBDriver,
		DBSource: config.DBSource,
	}
	model.Setup(dbConfig)
}
