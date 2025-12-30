package main

import (
	"context"

	spellchecker "github.com/f1monkey/spellchecker/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/ziedyousfi/axidev-io-go/keyboard"
)

// App struct holds the application state
type App struct {
	ctx         context.Context
	CurrentWord *CurrentWord
}

// NewApp creates a new App instance
func NewApp(checker *spellchecker.Spellchecker, sender *keyboard.Sender) *App {
	app := &App{}
	app.CurrentWord = NewCurrentWord(app, checker, sender)
	return app
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// UpdateDisplay sends display text to the frontend
func (a *App) UpdateDisplay(text string, state string) {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "updateText", map[string]string{
		"text":  text,
		"state": state,
	})
}
