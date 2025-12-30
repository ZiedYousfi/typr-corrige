package main

import (
	"fmt"
	"strings"
	"time"

	spellchecker "github.com/f1monkey/spellchecker/v3"
	"github.com/ziedyousfi/axidev-io-go/keyboard"
)

type CurrentWord struct {
	Word      string
	App       *App
	StartTime time.Time
	Checker   *spellchecker.Spellchecker
	Sender    *keyboard.Sender
}

func NewCurrentWord(app *App, checker *spellchecker.Spellchecker, sender *keyboard.Sender) *CurrentWord {
	return &CurrentWord{
		App:     app,
		Checker: checker,
		Sender:  sender,
	}
}

func (w *CurrentWord) Callback(event keyboard.KeyEvent) {
	if !event.IsPress() {
		return
	}

	r := event.Rune()
	if r == ' ' {
		if w.Word != "" {
			wordLower := strings.ToLower(w.Word)
			isCorrect := w.Checker.IsCorrect(wordLower)

			fmt.Printf("\n=== Word: %s ===\n", w.Word)
			fmt.Printf("Spelling: ")
			if isCorrect {
				fmt.Println("✓ CORRECT")
			} else {
				fmt.Println("✗ INCORRECT")
				result := w.Checker.Suggest(wordLower, 3)
				if len(result.Suggestions) > 0 {
					words := make([]string, len(result.Suggestions))
					for i, match := range result.Suggestions {
						words[i] = match.Value
					}
					fmt.Printf("Suggestions: %v\n", words)

					if w.Sender != nil {
						correction := result.Suggestions[0].Value
						fmt.Printf("Auto-correcting to: %s\n", correction)
						w.correctWord(correction)
					}
				}
			}
			fmt.Println()

			w.Word = ""
			w.StartTime = time.Time{}
		}
	} else if r != 0 {
		if w.Word == "" {
			w.StartTime = time.Now()
		}
		w.Word += string(r)
	}

	w.updateDisplay()
}

func (w *CurrentWord) updateDisplay() {
	if w.App == nil {
		return
	}

	var displayText string
	var state string

	if w.Word == "" {
		displayText = "Waiting..."
		state = "waiting"
	} else {
		wordLower := strings.ToLower(w.Word)
		if w.Checker.IsCorrect(wordLower) {
			displayText = w.Word + " ✓"
			state = "correct"
		} else {
			result := w.Checker.Suggest(wordLower, 1)
			if len(result.Suggestions) > 0 {
				displayText = fmt.Sprintf("%s → %s", w.Word, result.Suggestions[0].Value)
				state = "suggestion"
			} else {
				displayText = w.Word + " ?"
				state = "incorrect"
			}
		}
	}

	w.App.UpdateDisplay(displayText, state)
}

func (w *CurrentWord) correctWord(correction string) {
	if w.Sender == nil {
		return
	}

	time.Sleep(50 * time.Millisecond)

	leftKey := keyboard.StringToKey("Left")
	backspaceKey := keyboard.StringToKey("Backspace")

	if err := w.Sender.Combo(keyboard.ModAlt|keyboard.ModShift, leftKey); err != nil {
		fmt.Printf("Error selecting word: %v\n", err)
		return
	}

	time.Sleep(20 * time.Millisecond)

	if err := w.Sender.Tap(backspaceKey); err != nil {
		fmt.Printf("Error deleting word: %v\n", err)
		return
	}

	time.Sleep(20 * time.Millisecond)

	if err := w.Sender.TypeText(correction + " "); err != nil {
		fmt.Printf("Error typing correction: %v\n", err)
		return
	}

	w.Sender.Flush()
}
