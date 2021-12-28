package main

import (
	_ "embed"
	"fmt"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var icon []byte

func initTray() {
	systray.Register(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("RectangleWin")
	systray.SetTooltip("RectangleWin")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	mQuit.SetIcon(icon)
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("clicked Quit")
		systray.Quit()
	}()
	fmt.Println("tray ready")
}

func onExit() {
	fmt.Println("onExit invoked")
}
