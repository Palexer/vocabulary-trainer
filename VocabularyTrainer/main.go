package main

import (
	"fyne.io/fyne/app"

	"fyne.io/fyne"
)

type vocabulary struct {
	Title           string     `json:"Title"`
	Vocabulary      [][]string `json:"Vocabulary"`
	FirstLanguage   string
	SecondLanguage  string
	CurrentLanguage string
}

func main() {
	a := app.NewWithID("io.github.palexer")
	win := a.NewWindow("Vocabulary Trainer")
	win.SetIcon(resourceIconPng)
	win.Resize(fyne.NewSize(560, 450))
	trainerUI := &UI{mainWin: win, app: a}
	trainerUI.init()
	win.SetContent(trainerUI.loadMainUI())
	win.ShowAndRun()
}
