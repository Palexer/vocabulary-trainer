package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// SetupUISettings creates the settings dialog for the application
func SetupUISettings() {
	windowSettings := App.NewWindow("Settings")
	windowSettings.Resize(fyne.Size{
		Width:  600,
		Height: 440,
	})

	settingsLabel := widget.NewLabel("Settings")
	themeSelectorLabel := widget.NewLabel("Theme")
	themeSelector := widget.NewSelect([]string{"Light", "Dark"}, func(selectedTheme string) {
		switch selectedTheme {
		case "Light":
			App.Settings().SetTheme(theme.LightTheme())
		case "Dark":
			App.Settings().SetTheme(theme.DarkTheme())
		}

		App.Preferences().SetString("Theme", selectedTheme)
	})

	themeSelector.SetSelected(App.Preferences().StringWithFallback("Theme", "Dark"))

	windowSettings.SetContent(
		widget.NewVBox(
			settingsLabel,
			widget.NewHBox(
				themeSelectorLabel,
				themeSelector,
			),
		))
	windowSettings.Show()
}

// func loadSettings() {
// 	App.Settings().SetTheme(App.Preferences().)
// }
