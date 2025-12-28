package main

import (
	"fmt"
	"log"

	typrio "github.com/ziedyousfi/typr-io-go"
)

type CurrentWord struct {
	Word  string
}

func main() {
	fmt.Println("Listening for keyboard events... (Press Space to clear word)")

	cw := &CurrentWord{}

	listener, err := typrio.NewListener()
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	err = listener.Start(cw.cb)
	if err != nil {
		log.Fatal(err)
	}

	// Keep the program running
	select {}
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
		fmt.Printf("\rCurrent word: %-50s\n", w.Word)
	}
}