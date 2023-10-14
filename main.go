package main

import (
	"github.com/Hamedblue1381/hamed-url-shortener/api"
	"github.com/Hamedblue1381/hamed-url-shortener/db/model"
	"github.com/Hamedblue1381/hamed-url-shortener/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		panic("cannot load config:")
	}
	server, err := api.NewServer(config)
	if err != nil {
		panic("Failed to create the API server: " + err.Error())
	}
	dbConfig := model.Config{
		DBDriver: config.DBDriver,
		DBSource: config.DBSource,
	}
	model.Setup(dbConfig)

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		panic("Failed to start the server: " + err.Error())
	}
}
