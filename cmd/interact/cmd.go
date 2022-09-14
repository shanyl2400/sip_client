package main

import (
	"fmt"
	"os"
	"sipsimclient/service/devices"
	"sipsimclient/service/sip"
	"strings"

	"github.com/c-bata/go-prompt"
)

type Prompt struct {
	prompt        *prompt.Prompt
	sipClient     sip.SIPClient
	deviceManager devices.DeviceManager
}

var suggestions = []prompt.Suggest{
	// Command
	{Text: "exit", Description: "Exit the program"},
	{Text: "devices", Description: "Operate devices"},
	{Text: "http", Description: "Get data from database"},
}

func (p *Prompt) executor(in string) {
	in = strings.TrimSpace(in)

	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "devices":
		p.handleDevice(blocks[1:])
		return
	case "http":
		p.handleHTTP(blocks[1:])
	}
}
func (p *Prompt) completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}
func (p *Prompt) Run() {
	p.prompt.Run()
}

func NewPrompt() *Prompt {
	pt := &Prompt{
		deviceManager: devices.NewDeviceManager(),
		sipClient:     sip.NewSIPClient(),
	}
	p := prompt.New(
		pt.executor,
		pt.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("http-prompt"),
	)
	pt.deviceManager.Init()
	pt.prompt = p
	return pt
}
