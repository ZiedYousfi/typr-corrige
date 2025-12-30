package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	axidevio "github.com/ziedyousfi/axidev-io-go"
	"github.com/ziedyousfi/axidev-io-go/keyboard"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	axidevio.SetLogLevel(axidevio.LogLevelWarn)
	fmt.Println("Listening for keyboard events... (Press Space to clear word)")

	// Initialize spellchecker
	sc, err := NewFrenchSpellchecker()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d French words into dictionary\n", len(frenchWords))

	// Create sender for auto-correction
	sender, err := keyboard.NewSender()
	if err != nil {
		log.Fatal("Failed to create sender:", err)
	}
	defer sender.Close()

	// Request accessibility permissions if needed (macOS)
	caps := sender.Capabilities()
	if caps.NeedsAccessibilityPerm {
		fmt.Println("Requesting accessibility permissions...")
		if !sender.RequestPermissions() {
			fmt.Println("Warning: Permissions not granted, auto-correction may not work")
		}
	}

	// Create app instance
	app := NewApp(sc, sender)

	// Initialize keyboard listener
	listener, err := keyboard.NewListener()
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Start keyboard listener in goroutine
	go func() {
		err = listener.Start(app.CurrentWord.Callback)
		if err != nil {
			log.Printf("Listener error: %v", err)
		}
	}()

	// Create Wails application
	err = wails.Run(&options.App{
		Title:  "Typr Correct",
		Width:  400,
		Height: 100,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal("Error:", err)
	}
}
