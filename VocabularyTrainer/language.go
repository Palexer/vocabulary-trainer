package main

type language struct {
	Translation         string `json:"Translation"`
	Grammar             string `json:"Grammar"`
	Forward             string `json:"Forward"`
	Check               string `json:"Check"`
	SwitchLanguages     string `json:"SwitchLanguages"`
	OpenFile            string `json:"OpenFile"`
	VocabularyGenerator string `json:"VocabularyGenerator"`
	Random              string `json:"Random"`
	Correct             string `json:"Correct"`
	PartlyCorrect       string `json:"PartlyCorrect"`
	CorrectAnswer       string `json:"CorrectAnswer"`
	Wrong               string `json:"Wrong"`
	ConfirmDone         string `json:"ConfirmDone"`
	ConfirmEnd          string `json:"ConfirmEnd"`
	Restart             string `json:"Restart"`
	WrongWords          string `json:"WrongWords"`
	EverythingCorrect   string `json:"EverythingCorrect"`
	WrongAnswers        string `json:"WrongAnswer"`
	FinishedWords       string `json:"FinishedWords"`
	CorrectAnswers      string `json:"CorrectAnswers"`

	EVocGenMobile        string `json:"EVocGenMobile"`
	EAlreadyChecked      string `json:"EAlreadyChecked"`
	EEnterCheck          string `json:"EEnterCheck"`
	ECheckBeforeContinue string `json:"ECheckBeforeContinue"`
	ECancelled           string `json:"ECancelled"`
	ENoContent           string `json:"ENoContent"`
	EWrongFile           string `json:"EWrongFile"`
	EWrongVocabulary     string `json:"EWrongVocabulary"`
	E2Audio              string `json:"E2Audio"`
	EOpenToUse           string `json:"EOpenToUse"`

	Settings        string `json:"Settings"`
	Theme           string `json:"Theme"`
	Language        string `json:"Language"`
	MoreInfo        string `json:"MoreInfo"`
	RestartRequired string `json:"RestartRequired"`
	RestartInfo     string `json:"RestartInfo"`

	Save               string `json:"Save"`
	Title              string `json:"Title"`
	ForeignWord        string `json:"ForeignWord"`
	FirstLanguage      string `json:"FirstLanguage"`
	SecondLanguage     string `json:"SecondLanguage"`
	SaveWord           string `json:"SaveWord"`
	Clear              string `json:"Clear"`
	RemoveLastEntry    string `json:"RemoveLastEntry"`
	AvailableLanguages string `json:"AvailableLanguages"`
	ShowWords          string `json:"ShowWords"`
	LastWords          string `json:"LastWords"`

	EWordAndTranslation string `json:"EWordAndTranslation"`
	ENoVocabulary       string `json:"ENoVocabulary"`
	EEnterLanguage      string `json:"EEnterLanguage"`
	EJSONExt            string `json:"EJSONExt"`
}
