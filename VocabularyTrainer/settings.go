package main

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// loadUISettings creates the settings dialog for the application
func (u *UI) loadUISettings() {
	u.winSettings = u.app.NewWindow("Settings")
	u.winSettings.Resize(fyne.NewSize(360, 250))

	settingsLabel := widget.NewLabel("Settings")
	infoLabel := widget.NewLabel("v2.1 | License: GPLv3")

	// theme selector
	themeSelectorLabel := widget.NewLabel("Theme")
	themeSelector := widget.NewSelect([]string{"Light", "Dark"}, func(selectedTheme string) {
		switch selectedTheme {
		case "Light":
			u.app.Settings().SetTheme(theme.LightTheme())
		case "Dark":
			u.app.Settings().SetTheme(theme.DarkTheme())
		}

		u.app.Preferences().SetString("Theme", selectedTheme)
	})
	themeSelector.SetSelected(u.app.Preferences().StringWithFallback("Theme", "Dark"))

	// language selector
	langSelectorLabel := widget.NewLabel("Language")
	langSelector := widget.NewSelect([]string{"English", "German"}, func(selectedLanguage string) {
		if selectedLanguage != u.app.Preferences().String("Language") {
			dialog.ShowInformation("Restart required", "You need to restart the application\nin order to change the language", u.winSettings)
		}
		u.app.Preferences().SetString("Language", selectedLanguage)
	})
	langSelector.SetSelected(u.app.Preferences().StringWithFallback("Language", "English"))

	githubLink := widget.NewHyperlink("more information on Github", u.parseURL("https://github.com/Palexer/vocabulary-trainer"))

	u.winSettings.SetContent(
		widget.NewVBox(
			settingsLabel,
			widget.NewHBox(
				themeSelectorLabel,
				layout.NewSpacer(),
				themeSelector,
			),
			widget.NewHBox(
				langSelectorLabel,
				layout.NewSpacer(),
				langSelector,
			),
			layout.NewSpacer(),
			widget.NewHBox(
				infoLabel,
				layout.NewSpacer(),
				githubLink,
			),
		))
	u.winSettings.Show()
}

func (u *UI) parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}
	return link
}
