package main

import (
	"sipsimclient/cmd"
	"sipsimclient/config"
	"sipsimclient/repository"
)

func main() {
	config.Load()
	repository.Init()
	defer repository.Close()

	pt := cmd.NewPrompt()
	pt.Run()
}
