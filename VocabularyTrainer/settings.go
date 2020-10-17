package main

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
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

	// theme selector
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

	githubLink := widget.NewHyperlink("More information on Github", parseURL("https://github.com/Palexer/vocabulary-trainer"))

	windowSettings.SetContent(
		widget.NewVBox(
			settingsLabel,
			widget.NewHBox(
				themeSelectorLabel,
				themeSelector,
			),
			layout.NewSpacer(),
			githubLink,
		))
	windowSettings.Show()
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	return link
}
