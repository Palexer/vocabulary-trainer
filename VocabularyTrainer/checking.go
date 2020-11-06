package main

import "strings"

// CheckTranslation checks the translation entered by the user against the right answer from the .json file
func CheckTranslation(input, correctAnswer string) bool {
	if input == correctAnswer {
		return true
	}

	editedInput := []string{}

	for _, word := range strings.Split(input, ",") {
		word = strings.TrimSuffix(word, " ")
		word = strings.TrimPrefix(word, " ")
		editedInput = append(editedInput, word)
	}

	if strings.Join(editedInput, ",") == correctAnswer {
		return true
	}

	for _, answer := range strings.Split(correctAnswer, ",") {
		if answer == input {
			return true
		}
	}
	return false
}

// CheckGrammar checks the translation entered by the user against the right answer from the .json file
func CheckGrammar(input, correctAnswer string) bool {
	if correctAnswer == "" && input == "" {
		return true

	} else if input == correctAnswer {
		return true

	} else {
		editedInput := []string{}

		for _, word := range strings.Split(input, ",") {
			word = strings.TrimSuffix(word, " ")
			word = strings.TrimPrefix(word, " ")
			editedInput = append(editedInput, word)
		}

		if strings.Join(editedInput, ",") == correctAnswer {
			return true
		}
	}
	return false
}
