#!/bin/bash
fyne bundle resources/icon.png > bundled.go
fyne bundle -append lang/de.json >> bundled.go
fyne bundle -append lang/en.json >> bundled.go
