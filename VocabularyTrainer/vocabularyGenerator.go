package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"fyne.io/fyne/driver/desktop"

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

// loadUIGenerator builds up the UI for the vocabulary generator
func (u *UI) loadUIGenerator() {
	u.writeIndex = 0

	u.winGenerator = u.app.NewWindow("Vocabulary Generator")
	u.winGenerator.Resize(fyne.NewSize(460, 350))
	u.winGenerator.SetIcon(resourceIconPng)

	u.saveFileBtn = widget.NewButtonWithIcon("Save File", theme.DocumentSaveIcon(), u.saveFile)

	u.titleInput = widget.NewEntry()
	u.titleInput.SetPlaceHolder("Title")

	u.foreignWordInput = widget.NewEntry()
	u.correctTranslationInput = widget.NewEntry()
	u.correctGrammarInput = widget.NewEntry()

	u.foreignWordInput.SetPlaceHolder("Foreign Word")
	u.correctTranslationInput.SetPlaceHolder("Translation")
	u.correctGrammarInput.SetPlaceHolder("Grammar")

	saveWordBtn := widget.NewButtonWithIcon("Save Word", theme.MailForwardIcon(), u.saveWord)

	clearBtn := widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
		u.writeIndex = 0
		u.newJSONFile.Title = ""
		u.newJSONFile.Vocabulary = u.newJSONFile.Vocabulary[:0]
	})

	backBtn := widget.NewButtonWithIcon("Remove last entry", theme.NavigateBackIcon(), func() {
		if len(u.newJSONFile.Vocabulary) == 0 {
			dialog.ShowError(errors.New("the file doesn't contain any vocabulary"), u.winGenerator)
			return
		}
		u.writeIndex--
		u.newJSONFile.Vocabulary = u.newJSONFile.Vocabulary[:u.writeIndex]
	})

	// keyboard shortcuts
	// save file using ctrl+s
	u.winGenerator.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: desktop.ControlModifier,
	}, func(_ fyne.Shortcut) {
		u.saveFile()
	})

	// next word / save word using ctrl+n
	u.winGenerator.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: desktop.ControlModifier,
	}, func(_ fyne.Shortcut) {
		u.saveWord()
	})

	u.winGenerator.SetContent(
		widget.NewVBox(
			u.saveFileBtn,
			u.titleInput,
			layout.NewSpacer(),
			u.foreignWordInput,
			u.correctTranslationInput,
			u.correctGrammarInput,
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewHBox(
				backBtn,
				clearBtn,
				layout.NewSpacer(),
				saveWordBtn,
			),
		))
	u.winGenerator.Show()
}

func (u *UI) saveFile() {
	if len(u.newJSONFile.Vocabulary) == 0 {
		dialog.ShowError(errors.New("the file doesn't contain any vocabulary"), u.winGenerator)
		return
	}

	saveFileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, u.winGenerator)
			return
		}

		if err == nil && writer == nil {
			return
		}

		err = u.writeJSONFile(writer)
		if err != nil {
			dialog.ShowError(err, u.winGenerator)
		}

	}, u.winGenerator)

	saveFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	saveFileDialog.Show()
}

func (u *UI) saveWord() {
	if u.foreignWordInput.Text == "" || u.correctTranslationInput.Text == "" {
		dialog.ShowError(
			errors.New("please and enter at least a foreign word and the translation of it"),
			u.winGenerator)
		return
	}

	u.newJSONFile.Title = u.titleInput.Text
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
	u.newJSONFile.Vocabulary = append(
		u.newJSONFile.Vocabulary,
		vocabularyInputsFinished)

	u.writeIndex++
	u.foreignWordInput.SetText("")
	u.correctGrammarInput.SetText("")
	u.correctTranslationInput.SetText("")
}

func (u *UI) writeJSONFile(f fyne.URIWriteCloser) error {
	if f.URI().Extension() != ".json" {
		return errors.New("the vocabulary files needs the .json file extension")
	}

	if f == nil {
		return errors.New("cancelled")
	}

	encodedJSONFile, err := json.MarshalIndent(u.newJSONFile, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.URI().String()[7:], encodedJSONFile, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
