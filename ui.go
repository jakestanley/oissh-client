package main

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	textinput string
	prompt    *widgets.Paragraph
)

func initUi() {

	textinput = ""

	prompt = widgets.NewParagraph()
	prompt.Title = "Prompt"
	prompt.Text = ">"

	layout()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
}

func layout() {
	prompt.SetRect(0, 0, 80, 4)
}

func clearInputText() {
	textinput = ""
	appendInputText("")
}

func appendInputText(input string) {
	textinput += input
	prompt.Text = fmt.Sprintf("> %s", textinput)
}

func submit() {
	processInput(textinput)
}

func inputUi() bool {
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<Escape>":
				return false
			case "<Space>":
				appendInputText(" ")
			case "<Enter>":
				submit()
				clearInputText()
			default:
				appendInputText(e.ID)
			}
			break
		}
	}
	return true
}

func renderUi() {

	// TODO resizing
	// ui.ResizeEvent()

	for {
		ui.Render(prompt)
		time.Sleep(time.Millisecond / 10)
	}
}
