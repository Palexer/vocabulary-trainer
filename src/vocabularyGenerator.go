package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
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
func SetupUIVocabularyGenerator(parentApp fyne.App) {
	windowVocGenerator := parentApp.NewWindow("Vocabulary Generator")
	windowVocGenerator.Resize(fyne.Size{
		Width:  600,
		Height: 440,
	})

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

		saveFileDialog.Show()
	})

	titleInput := widget.NewEntry()
	titleInput.SetPlaceHolder("Title")

	correctForeignWordInput := widget.NewEntry()
	correctTranslationInput := widget.NewEntry()
	correctGrammarInput := widget.NewEntry()

	correctForeignWordInput.SetPlaceHolder("Foreign Word")
	correctTranslationInput.SetPlaceHolder("Translation")
	correctGrammarInput.SetPlaceHolder("Grammar")

	nextWordBtn := widget.NewButtonWithIcon("Next Word", theme.MailForwardIcon(), func() {
		newJSONFile.Title = titleInput.Text
		newJSONFile.Vocabulary[writeIndex][0] = correctForeignWordInput.Text
		newJSONFile.Vocabulary[writeIndex][1] = correctTranslationInput.Text
		newJSONFile.Vocabulary[writeIndex][2] = correctGrammarInput.Text
		writeIndex++

		// test if it works
		fmt.Println(newJSONFile)
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
			widget.NewHBox(
				correctForeignWordInput,
				layout.NewSpacer(),
				correctTranslationInput,
				layout.NewSpacer(),
				correctGrammarInput,
			),
			layout.NewSpacer(),
			widget.NewHBox(
				backBtn,
				layout.NewSpacer(),
				nextWordBtn,
			),
		))

	windowVocGenerator.Show()
}

func writeJSONFile(f fyne.URIWriteCloser) error {
	if f == nil {
		return errors.New("cancelled")
	}

	jsonFile, err := json.MarshalIndent(newJSONFile, "", " ")
	if err != nil {
		return err
	}

	ioutil.WriteFile(f.Name()+".json", jsonFile, os.ModePerm)
	return nil
}
