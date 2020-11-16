package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/app"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/storage"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	htgotts "github.com/hegedustibor/htgo-tts"
)

type vocabulary struct {
	Title           string     `json:"Title"`
	Vocabulary      [][]string `json:"Vocabulary"`
	FirstLanguage   string
	SecondLanguage  string
	CurrentLanguage string
}

// UI represents the main whole GUI
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

	writeIndex int

	// main UI
	app                fyne.App
	mainWin            fyne.Window
	title              *widget.Label
	foreignWord        *widget.Label
	result             *widget.Label
	correctCounter     *widget.Label
	finishedCounter    *widget.Label
	inputTranslation   *widget.Entry
	inputGrammar       *widget.Entry
	continueBtn        *widget.Button
	checkBtn           *widget.Button
	switchLanguagesBtn *widget.Button
	speakBtn           *widget.Button

	// generator UI
	winGenerator            fyne.Window
	titleInput              *widget.Entry
	foreignWordInput        *widget.Entry
	correctTranslationInput *widget.Entry
	correctGrammarInput     *widget.Entry
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

	u.openFileToUseProgram = true
	u.userHasTry = true

	// set ctrl to super modifier on darwin hosts
	if runtime.GOOS == "darwin" {
		u.modkey = desktop.SuperModifier
	} else {
		u.modkey = desktop.ControlModifier
	}
}

func (u *UI) loadMainUI() *widget.Box {
	u.loadPreferences()

	u.title = widget.NewLabel("")
	u.result = widget.NewLabel("")
	u.correctCounter = widget.NewLabel("")
	u.finishedCounter = widget.NewLabel("")

	u.foreignWord = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	u.inputTranslation = widget.NewEntry()
	u.inputGrammar = widget.NewEntry()
	u.inputTranslation.SetPlaceHolder(u.lang.Translation)
	u.inputGrammar.SetPlaceHolder(u.lang.Grammar)

	u.continueBtn = widget.NewButtonWithIcon(u.lang.Forward, theme.NavigateNextIcon(), u.continueFunc)
	u.checkBtn = widget.NewButtonWithIcon(u.lang.Check, theme.ConfirmIcon(), u.checkBtnFunc)

	u.switchLanguagesBtn = widget.NewButton(u.lang.SwitchLanguages, func() {
		if u.vocabularyFile.CurrentLanguage == u.vocabularyFile.FirstLanguage {
			u.vocabularyFile.CurrentLanguage = u.vocabularyFile.SecondLanguage
		} else {
			u.vocabularyFile.CurrentLanguage = u.vocabularyFile.FirstLanguage
		}
		if u.langIndex == 0 {
			u.langIndex = 1
		} else {
			u.langIndex = 0
		}
		u.foreignWord.SetText(u.vocabularyFile.Vocabulary[u.index][u.langIndex])
	})

	openButton := widget.NewButtonWithIcon(u.lang.OpenFile, theme.FolderOpenIcon(), u.openFileFunc)

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		u.loadUISettings()
	})

	openGeneratorBtn := widget.NewButtonWithIcon(u.lang.VocabularyGenerator, theme.FileApplicationIcon(), func() {
		if runtime.GOOS == "android" {
			dialog.ShowError(errors.New(u.lang.EVocGenMobile), u.mainWin)
			return
		}
		u.loadUIGenerator()
	})

	randomWordsCheck := widget.NewCheck(u.lang.Random, func(checked bool) {
		if checked == true {
			u.random = true

		} else {
			u.random = false
		}
	})

	u.speakBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		go u.speak()
	})

	// keyboard shortcuts
	// continue using ctrl+f
	u.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyF,
		Modifier: u.modkey,
	}, func(_ fyne.Shortcut) {
		u.continueFunc()
	})

	// check using ctrl+d
	u.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyD,
		Modifier: u.modkey,
	}, func(_ fyne.Shortcut) {
		u.checkBtnFunc()
	})

	// open generator using ctrl+g
	u.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyG,
		Modifier: u.modkey,
	}, func(_ fyne.Shortcut) {
		u.loadUIGenerator()
	})

	// open file using ctrl+o
	u.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyO,
		Modifier: u.modkey,
	}, func(_ fyne.Shortcut) {
		u.openFileFunc()
	})

	// close application using ctrl+q
	u.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyQ,
		Modifier: u.modkey,
	}, func(_ fyne.Shortcut) {
		u.mainWin.Close()
	})

	// disable all inputs + buttons as long as there is no file opened
	u.checkBtn.Disable()
	u.continueBtn.Disable()
	u.inputGrammar.Disable()
	u.inputTranslation.Disable()
	u.switchLanguagesBtn.Disable()
	u.speakBtn.Disable()

	// return the widgets in a VBox layout
	return widget.NewVBox(
		openButton,
		u.title,
		u.foreignWord,
		u.inputTranslation,
		u.inputGrammar,
		widget.NewHBox(
			u.checkBtn,
			u.continueBtn,
			u.speakBtn,
			u.result,
			layout.NewSpacer(),
			u.switchLanguagesBtn,
		),
		u.correctCounter,
		u.finishedCounter,
		layout.NewSpacer(),
		widget.NewHBox(
			settingsButton,
			randomWordsCheck,
			layout.NewSpacer(),
			openGeneratorBtn,
		),
	)
}

func (u *UI) loadPreferences() {
	// set correct theme
	switch u.app.Preferences().String("Theme") {
	case "Dark":
		u.app.Settings().SetTheme(theme.DarkTheme())
	case "Light":
		u.app.Settings().SetTheme(theme.LightTheme())
	}

	// set correct language
	switch u.app.Preferences().String("Language") {
	case "German":
		json.Unmarshal(resourceDeJson.Content(), &u.lang)
	case "English":
		json.Unmarshal(resourceEnJson.Content(), &u.lang)
	}
}

func (u *UI) checkBtnFunc() {
	if u.openFileToUseProgram == true {
		return
	}

	if u.userHasTry == false {
		dialog.ShowError(errors.New(u.lang.EAlreadyChecked), u.mainWin)
		return
	}

	if u.inputTranslation.Text == "" || u.inputGrammar.Text == "" && u.vocabularyFile.Vocabulary[u.index][2] != "" {
		dialog.ShowError(errors.New(u.lang.EEnterCheck), u.mainWin)
		return
	}

	var checkTranslation bool
	if u.langIndex == 0 {
		checkTranslation = CheckTranslation(u.inputTranslation.Text, u.vocabularyFile.Vocabulary[u.index][1])
	} else {
		checkTranslation = CheckTranslation(u.inputTranslation.Text, u.vocabularyFile.Vocabulary[u.index][0])
	}

	checkGrammar := CheckGrammar(u.inputGrammar.Text, u.vocabularyFile.Vocabulary[u.index][2])

	if checkTranslation && checkGrammar {
		u.result.SetText(u.lang.Correct)
		u.correct++

	} else if checkTranslation {
		u.result.SetText(u.lang.PartlyCorrect)
		u.wrongWordsList = append(u.wrongWordsList, u.vocabularyFile.Vocabulary[u.index])
		u.inputGrammar.SetText(u.lang.CorrectAnswer + u.vocabularyFile.Vocabulary[u.index][2])

	} else if checkGrammar {
		u.result.SetText(u.lang.PartlyCorrect)
		u.wrongWordsList = append(u.wrongWordsList, u.vocabularyFile.Vocabulary[u.index])
		u.inputTranslation.SetText(u.lang.CorrectAnswer + u.vocabularyFile.Vocabulary[u.index][1])

	} else {
		u.result.SetText(u.lang.Wrong)
		u.wrongWordsList = append(u.wrongWordsList, u.vocabularyFile.Vocabulary[u.index])
		u.inputTranslation.SetText(u.lang.CorrectAnswer + u.vocabularyFile.Vocabulary[u.index][1])
		u.inputGrammar.SetText(u.lang.CorrectAnswer + u.vocabularyFile.Vocabulary[u.index][2])
	}
	u.didCheck, u.userHasTry = true, false
}

func (u *UI) continueFunc() {
	if u.openFileToUseProgram == true {
		return
	}

	// done dialog
	if u.index+1 == len(u.vocabularyFile.Vocabulary) && u.random != true {

		// calculate the percentage of correct answers
		var percentage float64 = math.Round((float64(u.correct)/float64(u.finishedWords+1)*100.0)*100) / 100
		doneDialog := dialog.NewConfirm(
			u.lang.ConfirmDone, u.lang.ConfirmEnd+strconv.Itoa(u.correct)+"/"+strconv.Itoa(u.finishedWords+1)+" ("+(strconv.FormatFloat(percentage, 'f', -1, 64))+"%)"+u.lang.Restart,
			func(restart bool) {
				u.index, u.correct, u.finishedWords = 0, 0, 0
				u.correctCounter.SetText("")
				u.finishedCounter.SetText("")
				u.inputGrammar.SetText("")
				u.inputTranslation.SetText("")
				u.result.SetText("")

				if restart == true {
					u.correct, u.index, u.finishedWords = 0, 0, 0
					u.foreignWord.SetText(u.vocabularyFile.Vocabulary[u.index][u.langIndex])

				} else {
					u.foreignWord.SetText("")
					u.title.SetText("")
					u.openFileToUseProgram = true
				}

				// append wrong words to a string
				var wrongWords string
				for i := range u.wrongWordsList {
					wrongWords = wrongWords + "\n" + strings.Join(u.wrongWordsList[i], " - ")
				}

				if len(u.wrongWordsList) == 0 {
					dialog.NewInformation(u.lang.WrongWords, u.lang.EverythingCorrect, u.mainWin)
				} else {
					dialog.NewInformation(u.lang.WrongWords, u.lang.WrongAnswers+wrongWords, u.mainWin)
				}

			}, u.mainWin)

		doneDialog.Show()

	} else {
		// forward usually
		if u.didCheck == false {
			dialog.ShowError(errors.New(u.lang.ECheckBeforeContinue), u.mainWin)
			return
		}

		if u.random {
			u.index = rand.Intn(len(u.vocabularyFile.Vocabulary))
			u.foreignWord.SetText(u.vocabularyFile.Vocabulary[u.index][u.langIndex])

		} else {
			u.index++
			u.foreignWord.SetText(u.vocabularyFile.Vocabulary[u.index][u.langIndex])
		}

		u.finishedWords++
		// cleanup
		u.inputTranslation.SetText("")
		u.inputGrammar.SetText("")
		u.result.SetText("")
	}
	u.finishedCounter.SetText(u.lang.FinishedWords + strconv.Itoa(u.finishedWords) + "/" + strconv.Itoa(len(u.vocabularyFile.Vocabulary)))
	u.correctCounter.SetText(u.lang.CorrectAnswers + strconv.Itoa(u.correct) + "/" + strconv.Itoa(u.finishedWords))
	if u.random {
		u.finishedCounter.Hide()
	}
	u.didCheck, u.userHasTry = false, true
}

func (u *UI) openFileFunc() {
	openFileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader == nil {
			return
		}
		if err != nil {
			dialog.ShowError(err, u.mainWin)
			return
		}

		err = u.openFile(reader)
		if err != nil {
			dialog.ShowError(err, u.mainWin)
			return
		}

		// activate inputs + buttons when a file is opened; cleanup
		u.checkBtn.Enable()
		u.continueBtn.Enable()
		u.switchLanguagesBtn.Enable()
		u.speakBtn.Enable()
		u.inputGrammar.Enable()
		u.inputTranslation.Enable()
		u.inputGrammar.SetText("")
		u.inputTranslation.SetText("")
		u.correctCounter.SetText("")
		u.finishedCounter.SetText("")
		u.index, u.correct, u.finishedWords, u.langIndex = 0, 0, 0, 0
		u.openFileToUseProgram = false

		u.title.SetText(u.vocabularyFile.Title)
		u.foreignWord.SetText(u.vocabularyFile.Vocabulary[u.index][u.langIndex])

	}, u.mainWin)

	openFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	openFileDialog.Show()
}

func (u *UI) openFile(f fyne.URIReadCloser) error {
	if f == nil {
		return errors.New(u.lang.ECancelled)
	}

	byteData, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if byteData == nil {
		return errors.New(u.lang.ENoContent)
	}

	json.Unmarshal(byteData, &u.vocabularyFile)
	u.vocabularyFile.CurrentLanguage = u.vocabularyFile.FirstLanguage

	if len(u.vocabularyFile.Vocabulary) == 0 {
		return errors.New(u.lang.EWrongFile)
	}

	for i := 0; i < len(u.vocabularyFile.Vocabulary); i++ {
		if len(u.vocabularyFile.Vocabulary[i]) != 3 {
			return errors.New(u.lang.EWrongVocabulary + strconv.Itoa(i+1) + " )")
		}
	}
	return nil
}

func (u *UI) speak() {
	s := htgotts.Speech{Folder: os.TempDir(), Language: u.vocabularyFile.CurrentLanguage}
	s.Speak(u.foreignWord.Text)
	err := u.playAudio(os.TempDir() + "/" + u.foreignWord.Text + ".mp3")
	if err != nil {
		dialog.ShowError(err, u.mainWin)
	}
	os.Remove(os.TempDir() + "/" + u.foreignWord.Text + ".mp3")
	u.audioBusy = false
}

func (u *UI) playAudio(file string) error {
	if u.audioBusy {
		return errors.New(u.lang.E2Audio)
	}

	u.audioBusy = true
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	if u.didSpeakerInit == false {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		u.didSpeakerInit = true
	} else {
		speaker.Unlock()
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	defer speaker.Lock()
	return nil
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
