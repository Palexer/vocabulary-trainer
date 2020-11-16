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
	Title          string
	FirstLanguage  string
	SecondLanguage string
	Vocabulary     [][]string
}

// loadUIGenerator builds up the UI for the vocabulary generator
func (u *UI) loadUIGenerator() {
	u.writeIndex = 0

	u.winGenerator = u.app.NewWindow(u.lang.VocabularyGenerator)
	u.winGenerator.Resize(fyne.NewSize(510, 410))
	u.winGenerator.SetIcon(resourceIconPng)

	u.saveFileBtn = widget.NewButtonWithIcon(u.lang.Save, theme.DocumentSaveIcon(), u.saveFile)

	u.titleInput = widget.NewEntry()
	u.titleInput.SetPlaceHolder(u.lang.Title)

	u.foreignWordInput = widget.NewEntry()
	u.correctTranslationInput = widget.NewEntry()
	u.correctGrammarInput = widget.NewEntry()

	u.foreignWordInput.SetPlaceHolder(u.lang.ForeignWord)
	u.correctTranslationInput.SetPlaceHolder(u.lang.Translation)
	u.correctGrammarInput.SetPlaceHolder(u.lang.Grammar)

	u.langOneInput = widget.NewEntry()
	u.langTwoInput = widget.NewEntry()
	u.langOneInput.SetPlaceHolder(u.lang.FirstLanguage)
	u.langTwoInput.SetPlaceHolder(u.lang.SecondLanguage)

	saveWordBtn := widget.NewButtonWithIcon(u.lang.SaveWord, theme.MailForwardIcon(), u.saveWord)

	// clears all the previously entered vocabulary
	clearBtn := widget.NewButtonWithIcon(u.lang.Clear, theme.ContentClearIcon(), u.clear)

	// removes the last entered word
	backBtn := widget.NewButtonWithIcon(u.lang.RemoveLastEntry, theme.NavigateBackIcon(), func() {
		if len(u.newJSONFile.Vocabulary) == 0 {
			dialog.ShowError(errors.New(u.lang.ENoContent), u.winGenerator)
			return
		}
		u.writeIndex--
		u.newJSONFile.Vocabulary = u.newJSONFile.Vocabulary[:u.writeIndex]
	})

	availableLangLink := widget.NewHyperlink(u.lang.AvailableLanguages, u.parseURL("https://github.com/Palexer/vocabulary-trainer#available-languages-for-tts"))

	showWordsBtn := widget.NewButton(u.lang.ShowWords, func() {
		var enteredWords string
		for i := range u.newJSONFile.Vocabulary {
			if i > 10 {
				continue
			}
			enteredWords = enteredWords + "\n" + strings.Join(u.newJSONFile.Vocabulary[i], " - ")
		}

		dialog.ShowInformation(u.lang.LastWords, enteredWords, u.winGenerator)
	})

	// keyboard shortcuts
	// save file using ctrl+s
	u.winGenerator.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: u.modkey,
	}, func(_ fyne.Shortcut) {
		u.saveFile()
	})

	// next word / save word using ctrl+n
	u.winGenerator.Canvas().AddShortcut(&desktop.CustomShortcut{
		KeyName:  fyne.KeyN,
		Modifier: u.modkey,
	}, func(_ fyne.Shortcut) {
		u.saveWord()
	})

	// add the widgets to a VBox layout and set it as the content of the window
	u.winGenerator.SetContent(
		widget.NewVBox(
			u.saveFileBtn,
			u.titleInput,
			u.langOneInput,
			u.langTwoInput,
			availableLangLink,
			layout.NewSpacer(),
			u.foreignWordInput,
			u.correctTranslationInput,
			u.correctGrammarInput,
			layout.NewSpacer(),
			layout.NewSpacer(),
			widget.NewHBox(
				backBtn,
				clearBtn,
				showWordsBtn,
				layout.NewSpacer(),
				saveWordBtn,
			),
		))
	u.winGenerator.Show()

	// clear the vocabulary when the window gets closed
	u.winGenerator.SetOnClosed(u.clear)
}

func (u *UI) saveWord() {
	if u.foreignWordInput.Text == "" || u.correctTranslationInput.Text == "" {
		dialog.ShowError(
			errors.New(u.lang.EWordAndTranslation),
			u.winGenerator)
		return
	}

	// cut out the suffix of the languages if it exists
	u.langOneInput.Text = strings.TrimSuffix(u.langOneInput.Text, " ")
	u.langOneInput.Text = strings.TrimPrefix(u.langOneInput.Text, " ")
	u.langTwoInput.Text = strings.TrimSuffix(u.langTwoInput.Text, " ")
	u.langTwoInput.Text = strings.TrimPrefix(u.langTwoInput.Text, " ")

	vocabularyInputsFinished := []string{}

	// remove spaces if necessary
	for i, input := range []string{u.foreignWordInput.Text, u.correctTranslationInput.Text, u.correctGrammarInput.Text} {
		editedText := []string{}

		if i > 0 {
			for _, word := range strings.Split(input, ",") {
				word = strings.TrimPrefix(word, " ")
				word = strings.TrimSuffix(word, " ")
				editedText = append(editedText, word)
			}

		} else {
			input = strings.TrimPrefix(input, " ")
			input = strings.TrimSuffix(input, " ")
			editedText = []string{input}
		}
		vocabularyInputsFinished = append(vocabularyInputsFinished, strings.Join(editedText, ","))
	}

	// save languages + title in the type
	u.newJSONFile.FirstLanguage = u.langOneInput.Text
	u.newJSONFile.SecondLanguage = u.langTwoInput.Text
	u.newJSONFile.Title = u.titleInput.Text

	// append new vocabulary to the list in the type
	u.newJSONFile.Vocabulary = append(
		u.newJSONFile.Vocabulary,
		vocabularyInputsFinished)

	u.writeIndex++
	u.foreignWordInput.SetText("")
	u.correctGrammarInput.SetText("")
	u.correctTranslationInput.SetText("")
}

func (u *UI) saveFile() {
	if len(u.newJSONFile.Vocabulary) == 0 {
		dialog.ShowError(errors.New(u.lang.ENoVocabulary), u.winGenerator)
		return
	}

	if u.langOneInput.Text == "" || u.langTwoInput.Text == "" {
		dialog.ShowError(
			errors.New(u.lang.EEnterLanguage),
			u.winGenerator)
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

func (u *UI) writeJSONFile(f fyne.URIWriteCloser) error {
	if f == nil {
		return errors.New(u.lang.ECancelled)
	}

	if f.URI().Extension() != ".json" {
		os.Remove(f.URI().String()[7:])
		return errors.New(u.lang.EJSONExt)
	}

	encodedJSONFile, err := json.MarshalIndent(u.newJSONFile, "", " ") // encode + indent json
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.URI().String()[7:], encodedJSONFile, os.ModePerm) // write file
	if err != nil {
		return err
	}
	return nil
}

func (u *UI) clear() {
	u.writeIndex = 0
	u.newJSONFile.Title = ""
	u.newJSONFile.Vocabulary = u.newJSONFile.Vocabulary[:0]
}
