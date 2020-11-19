package main

import (
	"runtime"

	"fyne.io/fyne/app"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/widget"

	"fyne.io/fyne"
)

type enterEntry struct {
	widget.Entry
	onEnter func()
}

func newEnterEntry(enterFunc func()) *enterEntry {
	entry := &enterEntry{onEnter: enterFunc}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *enterEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	default:
		e.Entry.KeyDown(key)
	}

}

// UI represents the whole GUI
type UI struct {
	// vars
	vocabularyFile       vocabulary
	index                int
	finishedWords        int
	correct              int
	langIndex            int
	wrongWordsList       [][]string
	didCheck             bool
	openFileToUseProgram bool
	userHasTry           bool
	random               bool
	didSpeakerInit       bool
	audioBusy            bool
	modkey               desktop.Modifier
	check                bool

	writeIndex int

	// main UI
	app                fyne.App
	mainWin            fyne.Window
	title              *widget.Label
	foreignWord        *widget.Label
	result             *widget.Label
	correctCounter     *widget.Label
	finishedCounter    *widget.Label
	inputTranslation   *enterEntry
	inputGrammar       *enterEntry
	mainForwardBtn     *widget.Button // the button that switches between check and continue
	switchLanguagesBtn *widget.Button
	speakBtn           *widget.Button

	// generator UI
	winGenerator            fyne.Window
	titleInput              *widget.Entry
	foreignWordInput        *enterEntry
	correctTranslationInput *enterEntry
	correctGrammarInput     *enterEntry
	saveFileBtn             *widget.Button
	newJSONFile             jsonFile
	langOneInput            *widget.Entry
	langTwoInput            *widget.Entry

	// settings UI
	winSettings   fyne.Window
	themeSelector *widget.Select
	langSelector  *widget.Select

	// languages
	lang language
}

func (u *UI) init() {
	// variables
	u.index = 0
	u.finishedWords = 0
	u.correct = 0
	u.langIndex = 0
	u.didCheck = false
	u.openFileToUseProgram = false
	u.userHasTry = false
	u.random = false
	u.check = true

	u.openFileToUseProgram = true
	u.userHasTry = true

	// set ctrl to super modifier on darwin hosts
	if runtime.GOOS == "darwin" {
		u.modkey = desktop.SuperModifier
	} else {
		u.modkey = desktop.ControlModifier
	}
}

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
