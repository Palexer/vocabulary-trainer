# Vocabulary Trainer

## About

Vocabulary Trainer is a small application written in Go, that can help you with learning new vocabulary.
**Please note that the current state is not production ready.**

## Usage

### Formatting JSON files

You need to open a correctly formatted .json-file with the Vocabulary Trainer.
Here is an example of a correctly formatted .json-file:

```JSON
{
    "Title": "The title of the vocabulary",
    "Vocabulary": [
        ["here goes the foreign word", "and here the translation", "and here optional grammar"],
        ["word in foreign language", "the user should input this", "the user has to input this in the grammar field"],
        ["you can also use multiple options", "by,seperating,them,with,a,comma,like,this", ""],
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
- automatically display the correct answer if input was wrong
- keep buttons disabled if no file was opened