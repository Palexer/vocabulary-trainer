# Vocabulary Trainer

## About

Vocabulary Trainer is a small application written in Go, that can help you with learning new vocabulary.
**Please note that the current state is not production ready.**

## Installation

### Windows

- Download the exe file from the releases section and execute it.

### Linux

- Download the .tar.gz package from the releases section and extract it.

- After that, open a terminal in the folder, that contains the "Makefile"

- Execute the following command:
  
  ```bash
  sudo make install
  ```

### Compile it yourself

coming soon

## Usage

In order to use the Vocabulary Trainer you need a correctly formatted .json-file. You can see how to create one down below.
You then need to open the file with the Vocabulary Trainer and the program will show you the first foreign word. Enter a 
translation and, if you need to, also additional grammar. You can also enter multiple translations by separating them with a comma.
Depending on the file, you may or may not need to insert a space after the comma (at least for now). The grammar should be entered completely and exactly like in the correct answer. If you entered your translation and grammar, click on the check button to see if you're input was correct. If this is not the case, you will now see the correct answer.  Click on "Continue" to go to the next word. The Vocabulary Trainer will count your correct answers as well as your already finished words while you are practicing.

### Formatting JSON files

You need to open a correctly formatted .json-file with the Vocabulary Trainer.
Here is an example of a correctly formatted .json-file:

```JSON
{
    "Title": "The title of the vocabulary",
    "Vocabulary": [
        ["here goes the foreign word", "and here the translation", "and here optional grammar"],
        ["word in foreign language", "the user should input this", "the user has to input this in the grammar field"],
        ["you can also use multiple options", "by,separating,them,with,a,comma,like,this", ""],
    ]
}
```

_Note: You are currently NOT allowed to have spaces after the commas._
The last list entry is optional (additional grammar). If you don't use it for a word, 
you have to put an empty string there (like in the last example above).

## License

- GPL v3

## ToDo

- enable/disable buttons if they aren't useable or make them do nothing
- keep buttons disabled if no file was opened
- improve restart
- Icon
- one try -> disable buttons
- application for creating correctly formatted json-files
- choose random vocabulary
- keyboard shortcuts
- percentage at the end
- error if only two items for one word instead of three