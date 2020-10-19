package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type jsonFile struct {
	Title      string
	Vocabulary [][]string
}

var newJSONFile jsonFile
var writeIndex int

// SetupUIVocabularyGenerator builds up the UI for the vocabulary generator
func SetupUIVocabularyGenerator() {
	windowVocGenerator := App.NewWindow("Vocabulary Generator")
	windowVocGenerator.Resize(fyne.Size{
		Width:  600,
		Height: 440,
	})
	windowVocGenerator.SetIcon(resourceIconPng)

	saveFileBtn := widget.NewButtonWithIcon("Save File", theme.DocumentSaveIcon(), func() {
		saveFileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, windowVocGenerator)
				return
			}

			if err == nil && writer == nil {
				return
			}

			err = writeJSONFile(writer)
			if err != nil {
				dialog.ShowError(err, windowVocGenerator)
			}

		}, windowVocGenerator)

		saveFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		saveFileDialog.Show()
	})

	titleInput := widget.NewEntry()
	titleInput.SetPlaceHolder("Title")

	foreignWordInput := widget.NewEntry()
	correctTranslationInput := widget.NewEntry()
	correctGrammarInput := widget.NewEntry()

	foreignWordInput.SetPlaceHolder("Foreign Word")
	correctTranslationInput.SetPlaceHolder("Translation")
	correctGrammarInput.SetPlaceHolder("Grammar")

	saveWordBtn := widget.NewButtonWithIcon("Save Word", theme.MailForwardIcon(), func() {
		newJSONFile.Title = titleInput.Text

		vocabularyInputs := []string{foreignWordInput.Text, correctTranslationInput.Text, correctGrammarInput.Text}
		vocabularyInputsFinished := []string{}

		// remove spaces if necessary
		for i, input := range vocabularyInputs {
			editedText := []string{}

			if i > 0 {
				for _, word := range strings.Split(input, ",") {
					if strings.HasPrefix(word, " ") {
						word = strings.TrimPrefix(word, " ")
					}

					if strings.HasSuffix(word, " ") {
						word = strings.TrimSuffix(word, " ")
					}

					editedText = append(editedText, word)
				}

			} else {
				if strings.HasPrefix(input, " ") {
					input = strings.TrimPrefix(input, " ")
				}
				if strings.HasSuffix(input, " ") {
					input = strings.TrimSuffix(input, " ")
				}
				editedText = []string{input}
			}
			vocabularyInputsFinished = append(vocabularyInputsFinished, strings.Join(editedText, ","))
		}

		// append new vocabulary to struct
		newJSONFile.Vocabulary = append(
			newJSONFile.Vocabulary,
			vocabularyInputsFinished)

		writeIndex++
		foreignWordInput.SetText("")
		correctGrammarInput.SetText("")
		correctTranslationInput.SetText("")
	})

	backBtn := widget.NewButtonWithIcon("Remove last entry", theme.NavigateBackIcon(), func() {
		if writeIndex <= 0 {
			dialog.ShowError(errors.New("can't remove last word because there are no words yet"), windowVocGenerator)
			return
		}
		writeIndex--
		newJSONFile.Vocabulary = newJSONFile.Vocabulary[:writeIndex]
	})

	windowVocGenerator.SetContent(
		widget.NewVBox(
			saveFileBtn,
			titleInput,
			layout.NewSpacer(),
			foreignWordInput,
			correctTranslationInput,
			correctGrammarInput,
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewHBox(
				backBtn,
				layout.NewSpacer(),
				saveWordBtn,
			),
		))
	windowVocGenerator.Show()
}

func writeJSONFile(f fyne.URIWriteCloser) error {
	if f.URI().Extension() != ".json" {
		return errors.New("the vocabulary files needs the .json file extension")
	}

	if f == nil {
		return errors.New("cancelled")
	}

	encodedJSONFile, err := json.MarshalIndent(newJSONFile, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.URI().String()[7:], encodedJSONFile, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}