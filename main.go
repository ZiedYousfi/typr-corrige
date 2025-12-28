package main

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	typrio "github.com/ziedyousfi/typr-io-go"
)

type CurrentWord struct {
	Word string
	Text *canvas.Text
}

func main() {
	fmt.Println("Listening for keyboard events... (Press Space to clear word)")

	myApp := app.New()

	// Utilisation du driver desktop pour créer une fenêtre sans bordures
	var window fyne.Window
	if drv, ok := myApp.Driver().(desktop.Driver); ok {
		window = drv.CreateSplashWindow()
	} else {
		window = myApp.NewWindow("Prototypage")
	}

	window.SetFixedSize(true)
	window.SetPadded(false)

	text := canvas.NewText("Waiting...", color.White)
	text.TextSize = 24 // Un peu plus grand pour la lisibilité
	content := container.NewCenter(text)

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 100))

	// Comme Move() n'est pas dispo, CenterOnScreen est la seule
	// commande de positionnement garantie par Fyne.
	window.CenterOnScreen()

	cw := &CurrentWord{
		Text: text,
	}

	listener, err := typrio.NewListener()
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// On lance le listener dans une goroutine pour ne pas bloquer l'UI
	go func() {
		err = listener.Start(cw.cb)
		if err != nil {
			log.Printf("Listener error: %v", err)
		}
	}()

	window.ShowAndRun()
}

func (w *CurrentWord) cb(event typrio.KeyEvent) {
	if !event.IsPress() {
		return
	}

	r := event.Rune()
	if r == ' ' {
		fmt.Printf("\nSpace detected. Clearing word: %s\n", w.Word)
		w.Word = ""
	} else if r != 0 {
		w.Word += string(r)
	}

	if w.Text != nil {
		fyne.Do(func() {
			if w.Word == "" {
				w.Text.Text = "Waiting..."
			} else {
				w.Text.Text = w.Word
			}
			w.Text.Refresh()
		})
	}
}