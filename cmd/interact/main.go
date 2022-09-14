package main

import (
	"sipsimclient/config"
	"sipsimclient/repository"
)

func main() {
	config.Load()
	repository.Init()
	defer repository.Close()

	pt := NewPrompt()
	pt.Run()
}
