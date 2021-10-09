package main

import (
	"errors"
	"os"
	"time"

	"fyne.io/fyne/v2/dialog"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	htgotts "github.com/hegedustibor/htgo-tts"
)

func (u *UI) speak() {
	// create audio file
	s := htgotts.Speech{Folder: os.TempDir(), Language: u.vocabularyFile.CurrentLanguage}
	s.Speak(u.foreignWord.Text)

	// play audio file
	err := u.playAudio(os.TempDir() + "/" + u.foreignWord.Text + ".mp3")
	if err != nil {
		dialog.ShowError(err, u.mainWin)
	}

	// remove audio file
	defer os.Remove(os.TempDir() + "/" + u.foreignWord.Text + ".mp3")

	u.audioBusy = false
}

func (u *UI) playAudio(file string) error {
	if u.audioBusy {
		return errors.New(u.lang.E2Audio)
	}

	u.audioBusy = true
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	if u.didSpeakerInit == false {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		u.didSpeakerInit = true
	} else {
		speaker.Unlock()
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	defer speaker.Lock()
	return nil
}
