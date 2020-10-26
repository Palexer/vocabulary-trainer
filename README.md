# Vocabulary Trainer

## About

Vocabulary Trainer is a small application written in Go, that can help you with learning new vocabulary.

![screenshot](screenshot.png)

## Installation

The VocabularyTrainer is currently available for Windows, macOS, Linux and Android. 
Please note, that the focus is not on mobile, that's why the Android version isn't that polished.

### Windows

- Download the .exe.zip file from the [releases section](https://github.com/Palexer/vocabulary-trainer/releases), unzip it and execute the .exe file it.

### Linux

- Download the .tar.gz package from the [releases section](https://github.com/Palexer/vocabulary-trainer/releases) and extract it.

- After that, open a terminal in the folder, that contains the "Makefile"

- Execute the following command:
  
  ```bash
  sudo make install
  ```

### macOS

- Download the .app.zip file from the [releases section](https://github.com/Palexer/vocabulary-trainer/releases) and unzip it. You can now drag it into your applications folder and execute it. 
  _(Alternatively you can directly execute the file from anywhere)_

### Android

- Download the .apk file from the [releases section](https://github.com/Palexer/vocabulary-trainer/releases) and install it on your Android device. 

### Compile it yourself

1. Install the Go compiler.
2. Install the fyne GUI toolkit and its dependencies.
3. Clone this repository.
4. Run ```go build .``` in the VocabularyTrainer directory.
5. Run ```fyne package -os darwin -icon resources/icon.png``` 
    _Note: You can replace ```darwin``` with ```windows``` or ```linux``` to get the packages for those platforms. For cross compiling as well as compiling for Android, please take a look at the fyne.io documentation._

## Usage

In order to use the Vocabulary Trainer you need a correctly formatted .json-file. You can see how to create one down below.
You then need to open the file with the Vocabulary Trainer and the program will show you the first foreign word. Enter a translation and, if you need to, also additional grammar. You can also enter multiple translations by separating them with a comma.
The grammar has to be entered completely. If you entered your translation and grammar, click on the check button to see if you're input was correct. If this is not the case, you will now see the correct answer.  Click on "Continue" to go to the next word. The Vocabulary Trainer will count your correct answers as well as your already finished words while you are practicing.

### Creating JSON files

You need to open a correctly formatted .json-file with the Vocabulary Trainer.
**You can either use the Vocabulary Generator, that is built in to the Vocabulary Trainer, or manually create one like shown below.**
Here is an example of a correctly formatted .json-file:

```JSON
{
 "Title": "Test File",
 "Vocabulary": [
  [
   "foreign word 1",
   "translation 1",
   "grammar1"
  ],
  [
   "foreign word 2",
   "translation2,alternative translation2",
   "grammar2,alternative grammar2"
  ],
  [
  "foreign word 3",
  "translation3",
  ""
  ]
 ]
}
```

_Note: You are NOT allowed to have spaces after the commas in the .json-files._
The last list entry is optional (additional grammar). If you don't use it for a word, 
you have to put an empty string there (like in the last example above).

To avoid errors, please use the Vocabulary Generator.

## License

- GPL v3

## Help

If you have any problems with the program or the repository you can create an [issue](https://github.com/Palexer/vocabulary-trainer/issues) on the Github page.

## ToDo

- keyboard shortcuts
- German language support
