package main

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// loadUISettings creates the settings dialog for the application
func (u *UI) loadUISettings() {
	winSettings := App.NewWindow("Settings")
	winSettings.Resize(fyne.NewSize(600, 440))

	settingsLabel := widget.NewLabel("Settings")
	infoLabel := widget.NewLabel("v1.2 | License: GPLv3")

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

	winSettings.SetContent(
		widget.NewVBox(
			settingsLabel,
			widget.NewHBox(
				themeSelectorLabel,
				themeSelector,
			),
			layout.NewSpacer(),
			widget.NewHBox(
				infoLabel,
				layout.NewSpacer(),
				githubLink,
			),
		))
	winSettings.Show()
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	return link
}
