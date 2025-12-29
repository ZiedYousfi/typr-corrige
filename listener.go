package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	spellchecker "github.com/f1monkey/spellchecker/v3"
	typrio "github.com/ziedyousfi/typr-io-go"
)

// CurrentWord tracks the currently typed word and its metrics
type CurrentWord struct {
	Word      string
	Text      *canvas.Text
	StartTime time.Time
	Checker   *spellchecker.Spellchecker
}

// NewCurrentWord creates a new CurrentWord instance
func NewCurrentWord(text *canvas.Text, checker *spellchecker.Spellchecker) *CurrentWord {
	return &CurrentWord{
		Text:    text,
		Checker: checker,
	}
}

// Callback handles keyboard events for the typing application
func (w *CurrentWord) Callback(event typrio.KeyEvent) {
	if !event.IsPress() {
		return
	}

	r := event.Rune()
	if r == ' ' {
		// Word completed - calculate speed and check spelling
		if w.Word != "" {
			// Check spelling
			wordLower := strings.ToLower(w.Word)
			isCorrect := w.Checker.IsCorrect(wordLower)

			fmt.Printf("\n=== Word: %s ===\n", w.Word)
			fmt.Printf("Spelling: ")
			if isCorrect {
				fmt.Println("✓ CORRECT")
			} else {
				fmt.Println("✗ INCORRECT")
				// Get suggestions (with max 3 results)
				result := w.Checker.Suggest(wordLower, 3)
				if len(result.Suggestions) > 0 {
					words := make([]string, len(result.Suggestions))
					for i, match := range result.Suggestions {
						words[i] = match.Value
					}
					fmt.Printf("Suggestions: %v\n", words)
				}
			}
			fmt.Println()

			w.Word = ""
			w.StartTime = time.Time{} // Reset start time
		}
	} else if r != 0 {
		// Start timing on first character
		if w.Word == "" {
			w.StartTime = time.Now()
		}
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
