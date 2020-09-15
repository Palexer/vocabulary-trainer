package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

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
	vocab   vocabulary
	index   int
	correct int
)

func setupUI() {
	app := app.New()
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

	checkButton := widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), func() {
		checkTranslation := checkTranslation(inputTranslation.Text, vocab.Vocabulary[index][1])
		checkGrammar := checkGrammar(inputGrammar.Text, vocab.Vocabulary[index][2])

		if checkTranslation && checkGrammar {
			result.SetText("Correct")
			correct++

		} else if checkTranslation || checkGrammar {
			result.SetText("Partly correct")

		} else {
			result.SetText("Wrong")
		}
	})

	continueButton := widget.NewButtonWithIcon("Continue", theme.NavigateNextIcon(), func() {
		if index == len(vocab.Vocabulary[index]) {
			doneDialog := dialog.NewConfirm("Done.", "You reached the end of the vocabulary list. Restart?", func(restart bool) {
				index, correct = 0, 0

				correctCounter.SetText("")
				finishedCounter.SetText("")
				inputGrammar.SetText("")
				inputTranslation.SetText("")

				if restart == true {
					foreignWord.SetText(vocab.Vocabulary[index][0])
				} else {
					foreignWord.SetText("")
					// disable buttons
				}
			}, window)

			doneDialog.Show()

		} else {
			index++
			foreignWord.SetText(vocab.Vocabulary[index][0])

			finishedCounter.SetText("Finished words: " + strconv.Itoa(index) + "/" + strconv.Itoa(len(vocab.Vocabulary)))
			correctCounter.SetText("Correct answers: " + strconv.Itoa(correct) + "/" + strconv.Itoa(index+1))

			// cleanup
			inputTranslation.SetText("")
			inputGrammar.SetText("")
			result.SetText("")
		}
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

			fileOpened(reader)
			title.SetText(vocab.Title)
			foreignWord.SetText(vocab.Vocabulary[index][0])

		}, window)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fileDialog.Show()

		// activate inputs + buttons when a file is opened; cleanup
		checkButton.Enable()
		continueButton.Enable()
		inputGrammar.Enable()
		inputTranslation.Enable()
		inputGrammar.SetText("")
		inputTranslation.SetText("")
		correctCounter.SetText("")
		finishedCounter.SetText("")
		index, correct = 0, 0
	})

	// enable all inputs + buttons as long as there is no file opened
	checkButton.Disable()
	continueButton.Disable()
	inputGrammar.Disable()
	inputTranslation.Disable()

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
		))

	window.ShowAndRun()
}

func fileOpened(f fyne.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}

	byteData, err := ioutil.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load text data", err)
		return
	}
	if byteData == nil {
		return
	}

	json.Unmarshal(byteData, &vocab)
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
