package main

import "strings"

// CheckTranslation checks the translation entered by the user against the right answer from the .json file
func CheckTranslation(inp, correctAnswer string) bool {
	if inp == correctAnswer {
		return true
	}

	inpSplitted := strings.Split(inp, ",")

	for i := 0; i < len(inpSplitted); i++ {
		if string(inpSplitted[i][0]) == " " {
			inpSplitted[i] = inpSplitted[i][1:]
		}
	}

	if strings.Join(inpSplitted, ",") == correctAnswer {
		return true
	}

	for _, answer := range strings.Split(correctAnswer, ",") {
		if answer == inp {
			return true
		}
	}
	return false
}

// CheckGrammar checks the translation entered by the user against the right answer from the .json file
func CheckGrammar(inp, correctAnswer string) bool {
	if correctAnswer == "" && inp == "" {
		return true

	} else if inp == correctAnswer {
		return true

	} else {
		inpSplitted := strings.Split(inp, ",")

		for i := 0; i < len(inpSplitted); i++ {
			if string(inpSplitted[i][0]) == " " {
				inpSplitted[i] = inpSplitted[i][1:]
			}
		}

		if strings.Join(inpSplitted, ",") == correctAnswer {
			return true
		}
	}
	return false
}
