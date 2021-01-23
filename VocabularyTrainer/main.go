package main

import (
	"encoding/json"
	"runtime"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2"
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
	vocabularyFile vocabulary
	index          int64
	finishedWords  int
	correct        int
	langIndex      int
	wrongWordsList [][]string
	didCheck       bool
	userHasTry     bool
	random         bool
	didSpeakerInit bool
	audioBusy      bool
	modkey         desktop.Modifier
	check          bool

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
	separator          *widget.Separator
	randomWordsCheck   *widget.Check

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
	u.random = false
	u.check = true
	u.userHasTry = true

	// set ctrl to super modifier on darwin hosts
	if runtime.GOOS == "darwin" {
		u.modkey = desktop.SuperModifier
	} else {
		u.modkey = desktop.ControlModifier
	}
}

func (u *UI) loadPreferences() {
	// set correct theme
	switch u.app.Preferences().String("Theme") {
	case "Dark":
		u.app.Settings().SetTheme(theme.DarkTheme())
	case "Light":
		u.app.Settings().SetTheme(theme.LightTheme())
	default:
		u.app.Settings().SetTheme(theme.DarkTheme()) // default theme is dark
	}

	// set correct language
	switch u.app.Preferences().String("Language") {
	case "German":
		json.Unmarshal(resourceDeJson.Content(), &u.lang)
	case "English":
		json.Unmarshal(resourceEnJson.Content(), &u.lang)
	default:
		json.Unmarshal(resourceEnJson.Content(), &u.lang) // default language is English
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
	a := app.NewWithID("io.github.palexer.vocabulary-trainer")
	win := a.NewWindow("Vocabulary Trainer")
	win.SetIcon(resourceIconPng)
	win.Resize(fyne.NewSize(560, 450))
	trainerUI := &UI{mainWin: win, app: a}
	trainerUI.init()
	win.SetContent(trainerUI.loadMainUI())
	win.ShowAndRun()
}
