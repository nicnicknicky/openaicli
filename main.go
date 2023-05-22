package main

import (
	"playground/openaicli/chatgpt"
	"playground/openaicli/tui"
)

func main() {
	tui.OAIClient = chatgpt.NewOpenAIClient()

	if err := tui.StartTUI(); err != nil {
		panic(err)
	}
}
