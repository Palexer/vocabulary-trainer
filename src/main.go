package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type vocabulary struct {
	Title      string     `json:"Title"`
	Vocabulary [][]string `json:"Vocabulary"`
}

var (
	vocabularyFile vocabulary
	index          int
	correct        int
)

func setupUI() {
	app := app.New()
	// app.SetIcon()
	window := app.NewWindow("Vocabulary Trainer")
	window.Resize(fyne.Size{
		Width:  640,
		Height: 480})

	title := widget.NewLabel("")
	foreignWord := widget.NewLabel("")
	result := widget.NewLabel("")
	correctCounter := widget.NewLabel("")
	finishedCounter := widget.NewLabel("")

	inputTranslation := widget.NewEntry()
	inputTranslation.SetPlaceHolder("Translation")
	inputGrammar := widget.NewEntry()
	inputGrammar.SetPlaceHolder("Grammar")

	continueButton := widget.NewButtonWithIcon("Continue", theme.NavigateNextIcon(), func() {
		if index+1 == len(vocabularyFile.Vocabulary) {
			doneDialog := dialog.NewConfirm(
				"Done.", "You reached the end of the vocabulary list. \n Correct answers: "+strconv.Itoa(correct)+"/"+strconv.Itoa(index+1)+"\n Restart?",
				func(restart bool) {
					index, correct = 0, 0

					correctCounter.SetText("")
					finishedCounter.SetText("")
					inputGrammar.SetText("")
					inputTranslation.SetText("")

					if restart == true {
						foreignWord.SetText(vocabularyFile.Vocabulary[index][0])

						// make everything an empty string in vocabularyFile
						vocabularyFile.Title = ""
						for i := 0; i < len(vocabularyFile.Vocabulary); i++ {
							for j := 0; j < len(vocabularyFile.Vocabulary[i]); j++ {
								vocabularyFile.Vocabulary[i][j] = ""
							}
						}
					} else {
						foreignWord.SetText("")
						// disable buttons
					}
				}, window)

			doneDialog.Show()

		} else {
			index++
			foreignWord.SetText(vocabularyFile.Vocabulary[index][0])

			finishedCounter.SetText("Finished words: " + strconv.Itoa(index) + "/" + strconv.Itoa(len(vocabularyFile.Vocabulary)))
			correctCounter.SetText("Correct answers: " + strconv.Itoa(correct) + "/" + strconv.Itoa(index))

			// cleanup
			inputTranslation.SetText("")
			inputGrammar.SetText("")
			result.SetText("")
		}
	})

	checkButton := widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), func() {
		checkTranslation := checkTranslation(inputTranslation.Text, vocabularyFile.Vocabulary[index][1])
		checkGrammar := checkGrammar(inputGrammar.Text, vocabularyFile.Vocabulary[index][2])

		if checkTranslation && checkGrammar {
			result.SetText("Correct")
			correct++

		} else if checkTranslation {
			result.SetText("Partly correct")
			inputGrammar.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][2])

		} else if checkGrammar {
			result.SetText("Partly correct")
			inputTranslation.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][1])

		} else {
			result.SetText("Wrong")
			inputTranslation.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][1])
			inputGrammar.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][2])
		}

		continueButton.Enable()
	})

	openButton := widget.NewButtonWithIcon("Open File", theme.FolderOpenIcon(), func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			err = fileOpened(reader)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			// activate inputs + buttons when a file is opened; cleanup
			checkButton.Enable()
			inputGrammar.Enable()
			inputTranslation.Enable()
			inputGrammar.SetText("")
			inputTranslation.SetText("")
			correctCounter.SetText("")
			finishedCounter.SetText("")
			index, correct = 0, 0

			title.SetText(vocabularyFile.Title)
			foreignWord.SetText(vocabularyFile.Vocabulary[index][0])
		}, window)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fileDialog.Show()
	})

	// enable all inputs + buttons as long as there is no file opened
	checkButton.Disable()
	continueButton.Disable()
	inputGrammar.Disable()
	inputTranslation.Disable()

	lightThemeBtn := widget.NewButton("Light Theme", func() {
		app.Settings().SetTheme(theme.LightTheme())
	})

	darkThemeBtn := widget.NewButton("Dark Theme", func() {
		app.Settings().SetTheme(theme.DarkTheme())
	})

	window.SetContent(
		widget.NewVBox(
			openButton,
			title,
			foreignWord,
			inputTranslation,
			inputGrammar,
			widget.NewHBox(
				checkButton,
				continueButton,
				result,
			),
			correctCounter,
			finishedCounter,
			layout.NewSpacer(),
			widget.NewHBox(
				lightThemeBtn,
				darkThemeBtn,
				layout.NewSpacer(),
			),
		))

	window.ShowAndRun()
}

func fileOpened(f fyne.URIReadCloser) error {
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

func checkTranslation(inp, correctAnswers string) bool {
	if inp == correctAnswers {
		return true
	}

	for _, answer := range strings.Split(correctAnswers, ",") {
		if answer == inp {
			return true
		}
	}
	return false
}

func checkGrammar(inp, correctAnswer string) bool {
	if correctAnswer == "" && inp == "" {
		return true
	}
	if inp == correctAnswer {
		return true
	}
	return false
}

func main() {
	setupUI()
}
