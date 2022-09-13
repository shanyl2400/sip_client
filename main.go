package main

import (
	"sipsimclient/config"
	"sipsimclient/http"
	"sipsimclient/repository"
)

func main() {
	config.Load()
	repository.Init()
	defer repository.Close()

	// pt := cmd.NewPrompt()
	// pt.Run()
	http.Start()
}
