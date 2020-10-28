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

// loadUIGenerator builds up the UI for the vocabulary generator
func (u *UI) loadUIGenerator() {
	winGenerator := u.app.NewWindow("Vocabulary Generator")
	winGenerator.Resize(fyne.NewSize(600, 440))
	winGenerator.SetIcon(resourceIconPng)

	saveFileBtn := widget.NewButtonWithIcon("Save File", theme.DocumentSaveIcon(), func() {
		saveFileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, winGenerator)
				return
			}

			if err == nil && writer == nil {
				return
			}

			err = writeJSONFile(writer)
			if err != nil {
				dialog.ShowError(err, winGenerator)
			}

		}, winGenerator)

		saveFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		saveFileDialog.Show()
	})

	u.titleInput = widget.NewEntry()
	u.titleInput.SetPlaceHolder("Title")

	u.foreignWordInput = widget.NewEntry()
	u.correctTranslationInput = widget.NewEntry()
	u.correctGrammarInput = widget.NewEntry()

	u.foreignWordInput.SetPlaceHolder("Foreign Word")
	u.correctTranslationInput.SetPlaceHolder("Translation")
	u.correctGrammarInput.SetPlaceHolder("Grammar")

	saveWordBtn := widget.NewButtonWithIcon("Save Word", theme.MailForwardIcon(), u.saveWord)

	backBtn := widget.NewButtonWithIcon("Remove last entry", theme.NavigateBackIcon(), func() {
		if writeIndex <= 0 {
			dialog.ShowError(errors.New("can't remove last word because there are no words yet"), winGenerator)
			return
		}
		writeIndex--
		newJSONFile.Vocabulary = newJSONFile.Vocabulary[:writeIndex]
	})

	winGenerator.SetContent(
		widget.NewVBox(
			saveFileBtn,
			u.titleInput,
			layout.NewSpacer(),
			u.foreignWordInput,
			u.correctTranslationInput,
			u.correctGrammarInput,
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewHBox(
				backBtn,
				layout.NewSpacer(),
				saveWordBtn,
			),
		))
	winGenerator.Show()
}

func (u *UI) saveWord() {
	newJSONFile.Title = u.titleInput.Text

	vocabularyInputs := []string{u.foreignWordInput.Text, u.correctTranslationInput.Text, u.correctGrammarInput.Text}
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
	u.foreignWordInput.SetText("")
	u.correctGrammarInput.SetText("")
	u.correctTranslationInput.SetText("")
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
