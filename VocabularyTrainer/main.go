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

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type vocabulary struct {
	Title      string     `json:"Title"`
	Vocabulary [][]string `json:"Vocabulary"`
}

var (
	vocabularyFile       vocabulary
	index                int
	finishedWords        int
	correct              int
	wrongWordsList       [][]string
	didCheck             bool
	openFileToUseProgram bool = true
	userHasTry           bool = true
	random               bool
	// App is the main application, that contains all windows.
	App fyne.App = app.NewWithID("com.palexer.vocabularytrainer")
)

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

func main() {
	window := App.NewWindow("Vocabulary Trainer")
	window.SetIcon(resourceIconPng)
	window.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})

	// load settings
	// set correct theme
	switch App.Preferences().String("Theme") {
	case "Dark":
		App.Settings().SetTheme(theme.DarkTheme())
	case "Light":
		App.Settings().SetTheme(theme.LightTheme())
	}

	// create input fields and labels
	title := widget.NewLabel("")
	foreignWord := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	result := widget.NewLabel("")
	correctCounter := widget.NewLabel("")
	finishedCounter := widget.NewLabel("")

	inputTranslation := widget.NewEntry()
	inputTranslation.SetPlaceHolder("Translation")
	inputGrammar := widget.NewEntry()
	inputGrammar.SetPlaceHolder("Grammar")

	continueButton := widget.NewButtonWithIcon("Continue", theme.NavigateNextIcon(), func() {
		if openFileToUseProgram == true {
			return
		}

		// done dialog
		if index+1 == len(vocabularyFile.Vocabulary) && random != true {

			// calculate the percentage of correct answers
			var percentage float64 = math.Round((float64(correct)/float64(len(vocabularyFile.Vocabulary))*100.0)*100) / 100
			doneDialog := dialog.NewConfirm(
				"Done.", "You reached the end of the vocabulary list. \n Correct answers: "+strconv.Itoa(correct)+"/"+strconv.Itoa(finishedWords+1)+"("+(strconv.FormatFloat(percentage, 'f', -1, 64))+"%)"+"\n Restart?",
				func(restart bool) {
					index, correct, finishedWords = 0, 0, 0
					correctCounter.SetText("")
					finishedCounter.SetText("")
					inputGrammar.SetText("")
					inputTranslation.SetText("")
					result.SetText("")

					if restart == true {
						correct, index, finishedWords = 0, 0, 0
						foreignWord.SetText(vocabularyFile.Vocabulary[index][0])

					} else {
						foreignWord.SetText("")
						title.SetText("")
						openFileToUseProgram = true
					}

					// append wrong words to a string
					var wrongWords string
					for i := range wrongWordsList {
						wrongWords = wrongWords + "\n" + strings.Join(wrongWordsList[i], " - ")
					}

					if len(wrongWordsList) == 0 {
						dialog.NewInformation("Wrong Words", "You entered everything correctly.", window)
					} else {
						dialog.NewInformation("Wrong Words", "You didn't know the solution to the following words:\n"+wrongWords, window)
					}

				}, window)

			doneDialog.Show()

		} else {
			// forward usually
			if didCheck == false {
				dialog.ShowError(errors.New("please check your input before you continue"), window)
				return
			}

			if random {
				index = rand.Intn(len(vocabularyFile.Vocabulary))
				foreignWord.SetText(vocabularyFile.Vocabulary[index][0])

			} else {
				index++
				foreignWord.SetText(vocabularyFile.Vocabulary[index][0])
			}

			finishedWords++
			// cleanup
			inputTranslation.SetText("")
			inputGrammar.SetText("")
			result.SetText("")
		}
		finishedCounter.SetText("Finished words: " + strconv.Itoa(finishedWords) + "/" + strconv.Itoa(len(vocabularyFile.Vocabulary)))
		correctCounter.SetText("Correct answers: " + strconv.Itoa(correct) + "/" + strconv.Itoa(finishedWords))
		if random {
			finishedCounter.Hide()
		}
		didCheck, userHasTry = false, true
	})

	checkButton := widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), func() {
		if openFileToUseProgram == true {
			return
		}

		if userHasTry == false {
			dialog.ShowError(errors.New("you already checked your input"), window)
			return
		}

		if inputTranslation.Text == "" || inputGrammar.Text == "" && vocabularyFile.Vocabulary[index][2] != "" {
			dialog.ShowError(errors.New("please enter a translation / the grammar first"), window)
			return
		}

		checkTranslation := CheckTranslation(inputTranslation.Text, vocabularyFile.Vocabulary[index][1])
		checkGrammar := CheckGrammar(inputGrammar.Text, vocabularyFile.Vocabulary[index][2])

		if checkTranslation && checkGrammar {
			result.SetText("Correct")
			correct++

		} else if checkTranslation {
			result.SetText("Partly correct")
			wrongWordsList = append(wrongWordsList, vocabularyFile.Vocabulary[index])
			inputGrammar.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][2])

		} else if checkGrammar {
			result.SetText("Partly correct")
			wrongWordsList = append(wrongWordsList, vocabularyFile.Vocabulary[index])
			inputTranslation.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][1])

		} else {
			result.SetText("Wrong")
			wrongWordsList = append(wrongWordsList, vocabularyFile.Vocabulary[index])
			inputTranslation.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][1])
			inputGrammar.SetText("Correct answer: " + vocabularyFile.Vocabulary[index][2])
		}
		didCheck, userHasTry = true, false
	})

	openButton := widget.NewButtonWithIcon("Open File", theme.FolderOpenIcon(), func() {
		openFileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
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
			continueButton.Enable()
			inputGrammar.Enable()
			inputTranslation.Enable()
			inputGrammar.SetText("")
			inputTranslation.SetText("")
			correctCounter.SetText("")
			finishedCounter.SetText("")
			index, correct, finishedWords = 0, 0, 0
			openFileToUseProgram = false

			title.SetText(vocabularyFile.Title)

			foreignWord.SetText(vocabularyFile.Vocabulary[index][0])

		}, window)

		openFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		openFileDialog.Show()
	})

	// enable all inputs + buttons as long as there is no file opened
	checkButton.Disable()
	continueButton.Disable()
	inputGrammar.Disable()
	inputTranslation.Disable()

	settingsButton := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
		SetupUISettings()
	})

	openGeneratorBtn := widget.NewButtonWithIcon("Vocabulary Generator", theme.FileApplicationIcon(), func() {
		if runtime.GOOS == "android" {
			dialog.ShowError(errors.New("the vocabulary generator\n is not support on mobile\n operating systems"), window)
			return
		}
		SetupUIVocabularyGenerator()
	})

	randomWordsBtn := widget.NewCheck("Random Words (infinite)", func(checked bool) {
		if checked == true {
			random = true

		} else {
			random = false
		}
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
				settingsButton,
				randomWordsBtn,
				layout.NewSpacer(),
				openGeneratorBtn,
			),
		))
	window.ShowAndRun()
}
