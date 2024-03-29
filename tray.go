// Copyright 2022 Ahmet Alp Balkan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "embed"
	"fmt"

	"github.com/getlantern/systray"
	"github.com/gonutz/w32/v2"
)

//go:embed assets/tray_icon.ico
var icon []byte

const repo = "https://github.com/ahmetb/RectangleWin"

func initTray() {
	systray.Register(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("RectangleWin")
	systray.SetTooltip("RectangleWin")

	autorun, err := AutoRunEnabled()
	if err != nil {
		panic(err)
	}

	mRepo := systray.AddMenuItem("Documentation", "")
	go func() {
		for range mRepo.ClickedCh {
			if err := w32.ShellExecute(0, "open", repo, "", "", w32.SW_SHOWNORMAL); err != nil {
				fmt.Printf("failed to launch browser: (%d), %v\n", w32.GetLastError(), err)
			}
		}
	}()

	systray.AddSeparator()

	mAutoRun := systray.AddMenuItemCheckbox("Run on startup", "", autorun)
	go func() {
		for range mAutoRun.ClickedCh {
			if mAutoRun.Checked() {
				if err := AutoRunDisable(); err != nil {
					mAutoRun.SetTitle(err.Error())
					fmt.Printf("warn: autorun disable: %v\n", err)
					continue
				}
				fmt.Println("disabled autorun")
				mAutoRun.Uncheck()
			} else {
				if err := AutoRunEnable(); err != nil {
					mAutoRun.SetTitle(err.Error())
					fmt.Printf("warn: autorun enable: %v\n", err)
					continue
				}
				fmt.Println("enabled autorun")
				mAutoRun.Check()
			}

		}
	}()

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "")
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
