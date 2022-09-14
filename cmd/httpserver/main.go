package main

import (
	"sipsimclient/api"
	"sipsimclient/config"
	"sipsimclient/repository"
)

func main() {
	config.Load()
	repository.Init()
	defer repository.Close()

	api.Start()
}
