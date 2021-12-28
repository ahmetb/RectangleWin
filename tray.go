package main

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var icon []byte

func initTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("RectangleWin")
	systray.SetTooltip("RectangleWin")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	mQuit.SetIcon(icon)
}

func onExit() {
}
