package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// loadUISettings creates the settings dialog for the application
func (u *UI) loadUISettings() {
	u.winSettings = u.app.NewWindow(u.lang.Settings)
	u.winSettings.SetIcon(resourceIconPng)
	u.winSettings.Resize(fyne.NewSize(360, 250))

	settingsLabel := widget.NewLabel(u.lang.Settings)
	infoLabel := widget.NewLabel("v1.3.1 | License: BSD 3-Clause")

	// theme selector
	themeSelectorLabel := widget.NewLabel(u.lang.Theme)
	u.themeSelector = widget.NewSelect([]string{"System Default", "Light", "Dark"}, func(selectedTheme string) {
		switch selectedTheme {
		case "Light":
			u.app.Settings().SetTheme(theme.LightTheme())
		case "Dark":
			u.app.Settings().SetTheme(theme.DarkTheme())
		case "System Default":
			u.app.Settings().SetTheme(theme.DefaultTheme())
		}

		u.app.Preferences().SetString("Theme", selectedTheme)
	})
	u.themeSelector.SetSelected(u.app.Preferences().StringWithFallback("Theme", "System Default"))

	// language selector
	langSelectorLabel := widget.NewLabel(u.lang.Language)
	u.langSelector = widget.NewSelect([]string{"English", "German"}, func(selectedLanguage string) {
		if selectedLanguage != u.app.Preferences().String("Language") {
			dialog.ShowInformation(u.lang.RestartRequired, u.lang.RestartInfo, u.winSettings)
		}
		u.app.Preferences().SetString("Language", selectedLanguage)
	})
	u.langSelector.SetSelected(u.app.Preferences().StringWithFallback("Language", "English"))

	githubLink := widget.NewHyperlink(u.lang.MoreInfo, u.parseURL("https://github.com/Palexer/vocabulary-trainer"))

	u.winSettings.SetContent(
		container.NewVBox(
			settingsLabel,
			container.NewHBox(
				themeSelectorLabel,
				layout.NewSpacer(),
				u.themeSelector,
			),
			container.NewHBox(
				langSelectorLabel,
				layout.NewSpacer(),
				u.langSelector,
			),
			layout.NewSpacer(),
			container.NewHBox(
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
