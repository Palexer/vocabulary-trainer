package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/storage"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type vocabulary struct {
	Title      string     `json:"Title"`
	Vocabulary [][]string `json:"Vocabulary"`
}

// UI represents the main whole GUI
type UI struct {
	// main UI
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

	// generator UI
	titleInput              *widget.Entry
	foreignWordInput        *widget.Entry
	correctTranslationInput *widget.Entry
	correctGrammarInput     *widget.Entry
}

var (
	vocabularyFile       vocabulary
	index                int
	finishedWords        int
	correct              int
	langIndex            int
	wrongWordsList       [][]string
	didCheck             bool
	openFileToUseProgram bool = true
	userHasTry           bool = true
	random               bool

	// App is the main application, that contains all windows.
	App fyne.App = app.NewWithID("com.palexer.vocabularytrainer")
)

func (u *UI) loadMainUI() *widget.Box {
	u.loadPreferences()

	u.title = widget.NewLabel("")
	u.result = widget.NewLabel("")
	u.correctCounter = widget.NewLabel("")
	u.finishedCounter = widget.NewLabel("")

	u.foreignWord = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	u.inputTranslation = widget.NewEntry()
	u.inputGrammar = widget.NewEntry()
	u.inputTranslation.SetPlaceHolder("Translation")
	u.inputGrammar.SetPlaceHolder("Grammar")

	u.continueBtn = widget.NewButtonWithIcon("Continue", theme.NavigateNextIcon(), u.continueFunc)
	u.checkBtn = widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), u.checkBtnFunc)

	u.switchLanguagesBtn = widget.NewButton("Switch Languages", func() {
		if langIndex == 0 {
			langIndex = 1
		} else {
			langIndex = 0
		}
		u.foreignWord.SetText(vocabularyFile.Vocabulary[index][langIndex])
	})

	openButton := widget.NewButtonWithIcon("Open File", theme.FolderOpenIcon(), u.openFileFunc)

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		u.loadUISettings()
	})

	openGeneratorBtn := widget.NewButtonWithIcon("Vocabulary Generator", theme.FileApplicationIcon(), func() {
		if runtime.GOOS == "android" {
			dialog.ShowError(errors.New("the vocabulary generator\n is not support on mobile\n operating systems"), u.mainWin)
			return
		}
		u.loadUIVocabularyGenerator()
	})

	randomWordsCheck := widget.NewCheck("Random Words (infinite)", func(checked bool) {
		if checked == true {
			random = true

		} else {
			random = false
		}
	})

	// keyboard shortcuts
	// u.mainWin.Canvas().AddShortcut(&desktop.CustomShortcut{
	// 	KeyName:  fyne.KeyEnter,
	// 	Modifier: desktop.ControlModifier,
	// }, func(_ fyne.Shortcut) {
	// 	fmt.Println("enter")
	// 	u.checkBtn.OnTapped()
	// })

	// enable all inputs + buttons as long as there is no file opened
	u.checkBtn.Disable()
	u.continueBtn.Disable()
	u.inputGrammar.Disable()
	u.inputTranslation.Disable()
	u.switchLanguagesBtn.Disable()

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
	switch App.Preferences().String("Theme") {
	case "Dark":
		App.Settings().SetTheme(theme.DarkTheme())
	case "Light":
		App.Settings().SetTheme(theme.LightTheme())
	}
}

func (u *UI) checkBtnFunc() {
	if openFileToUseProgram == true {
		return
	}

	if userHasTry == false {
		dialog.ShowError(errors.New("you already checked your input"), u.mainWin)
		return
	}

	if u.inputTranslation.Text == "" || u.inputGrammar.Text == "" && vocabularyFile.Vocabulary[index][2] != "" {
		dialog.ShowError(errors.New("please enter a translation / the grammar first"), u.mainWin)
		return
	}

	var checkTranslation bool
	if langIndex == 0 {
		checkTranslation = CheckTranslation(u.inputTranslation.Text, vocabularyFile.Vocabulary[index][1])
	} else {
		checkTranslation = CheckTranslation(u.inputTranslation.Text, vocabularyFile.Vocabulary[index][0])
	}

	checkGrammar := CheckGrammar(u.inputGrammar.Text, vocabularyFile.Vocabulary[index][2])

	if checkTranslation && checkGrammar {
		u.result.SetText("Correct")
		correct++

	} else if checkTranslation {
		u.result.SetText("Partly correct")
		wrongWordsList = append(wrongWordsList, vocabularyFile.Vocabulary[index])
		u.inputGrammar.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][2])

	} else if checkGrammar {
		u.result.SetText("Partly correct")
		wrongWordsList = append(wrongWordsList, vocabularyFile.Vocabulary[index])
		u.inputTranslation.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][1])

	} else {
		u.result.SetText("Wrong")
		wrongWordsList = append(wrongWordsList, vocabularyFile.Vocabulary[index])
		u.inputTranslation.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][1])
		u.inputGrammar.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][2])
	}
	didCheck, userHasTry = true, false
}

func (u *UI) continueFunc() {
	if openFileToUseProgram == true {
		return
	}

	// done dialog
	if index+1 == len(vocabularyFile.Vocabulary) && random != true {

		// calculate the percentage of correct answers
		var percentage float64 = math.Round((float64(correct)/float64(finishedWords+1)*100.0)*100) / 100
		doneDialog := dialog.NewConfirm(
			"Done.", "You reached the end of the vocabulary list. \n Correct answers: "+strconv.Itoa(correct)+"/"+strconv.Itoa(finishedWords+1)+" ("+(strconv.FormatFloat(percentage, 'f', -1, 64))+"%)"+"\n Restart?",
			func(restart bool) {
				index, correct, finishedWords = 0, 0, 0
				u.correctCounter.SetText("")
				u.finishedCounter.SetText("")
				u.inputGrammar.SetText("")
				u.inputTranslation.SetText("")
				u.result.SetText("")

				if restart == true {
					correct, index, finishedWords = 0, 0, 0
					u.foreignWord.SetText(vocabularyFile.Vocabulary[index][langIndex])

				} else {
					u.foreignWord.SetText("")
					u.title.SetText("")
					openFileToUseProgram = true
				}

				// append wrong words to a string
				var wrongWords string
				for i := range wrongWordsList {
					wrongWords = wrongWords + "\n" + strings.Join(wrongWordsList[i], " - ")
				}

				if len(wrongWordsList) == 0 {
					dialog.NewInformation("Wrong Words", "You entered everything correctly.", u.mainWin)
				} else {
					dialog.NewInformation("Wrong Words", "You didn't know the solution to the following words:\n"+wrongWords, u.mainWin)
				}

			}, u.mainWin)

		doneDialog.Show()

	} else {
		// forward usually
		if didCheck == false {
			dialog.ShowError(errors.New("please check your input before you continue"), u.mainWin)
			return
		}

		if random {
			index = rand.Intn(len(vocabularyFile.Vocabulary))
			u.foreignWord.SetText(vocabularyFile.Vocabulary[index][langIndex])

		} else {
			index++
			u.foreignWord.SetText(vocabularyFile.Vocabulary[index][langIndex])
		}

		finishedWords++
		// cleanup
		u.inputTranslation.SetText("")
		u.inputGrammar.SetText("")
		u.result.SetText("")
	}
	u.finishedCounter.SetText("Finished words: " + strconv.Itoa(finishedWords) + "/" + strconv.Itoa(len(vocabularyFile.Vocabulary)))
	u.correctCounter.SetText("Correct answers: " + strconv.Itoa(correct) + "/" + strconv.Itoa(finishedWords))
	if random {
		u.finishedCounter.Hide()
	}
	didCheck, userHasTry = false, true
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

		err = u.fileOpened(reader)
		if err != nil {
			dialog.ShowError(err, u.mainWin)
			return
		}

		// activate inputs + buttons when a file is opened; cleanup
		u.checkBtn.Enable()
		u.continueBtn.Enable()
		u.switchLanguagesBtn.Enable()
		u.inputGrammar.Enable()
		u.inputTranslation.Enable()
		u.inputGrammar.SetText("")
		u.inputTranslation.SetText("")
		u.correctCounter.SetText("")
		u.finishedCounter.SetText("")
		index, correct, finishedWords = 0, 0, 0
		openFileToUseProgram = false

		u.title.SetText(vocabularyFile.Title)

		u.foreignWord.SetText(vocabularyFile.Vocabulary[index][langIndex])

	}, u.mainWin)

	openFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	openFileDialog.Show()
}

func (u *UI) fileOpened(f fyne.URIReadCloser) error {
	if f == nil {
		return errors.New("cancelled")
	}

	byteData, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if byteData == nil {
		return errors.New("the file does not have any content")
	}

	json.Unmarshal(byteData, &vocabularyFile)

	if len(vocabularyFile.Vocabulary) == 0 {
		return errors.New("the file does not contain any vocabulary or is not correctly formatted")
	}

	for i := 0; i < len(vocabularyFile.Vocabulary); i++ {
		if len(vocabularyFile.Vocabulary[i]) != 3 {
			return errors.New("the file contains vocabulary with too many or too less arguments (error in list item " + strconv.Itoa(i+1) + " )")
		}
	}
	return nil
}

func main() {
	win := App.NewWindow("Vocabulary Trainer")
	win.SetIcon(resourceIconPng)
	win.Resize(fyne.NewSize(800, 600))
	trainerUI := &UI{mainWin: win}
	win.SetContent(trainerUI.loadMainUI())
	win.ShowAndRun()
}
