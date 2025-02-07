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

// TODO make it possible to "go generate" on Windows (https://github.com/josephspurrier/goversioninfo/issues/52).
//go:generate /bin/bash -c "go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest -arm -64 -icon=assets/icon.ico - <<< '{}'"

package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/gonutz/w32/v2"

	"github.com/ahmetb/RectangleWin/w32ex"
)

var lastResized w32.HWND

func main() {
	runtime.LockOSThread() // since we bind hotkeys etc that need to dispatch their message here
	if !w32ex.SetProcessDPIAware() {
		panic("failed to set DPI aware")
	}

	autorun, err := AutoRunEnabled()
	if err != nil {
		panic(err)
	}
	fmt.Printf("autorun enabled=%v\n", autorun)
	printMonitors()

	edgeFuncs := [][]resizeFunc{
		{leftHalf, leftTwoThirds, leftOneThirds},
		{rightHalf, rightTwoThirds, rightOneThirds},
		{topHalf, topTwoThirds, topOneThirds},
		{bottomHalf, bottomTwoThirds, bottomOneThirds}}
	// edgeFuncTurn := make([]int, len(edgeFuncs))
	// cornerFuncs := [][]resizeFunc{
	// 	{topLeftHalf, topLeftTwoThirds, topLeftOneThirds},
	// 	{topRightHalf, topRightTwoThirds, topRightOneThirds},
	// 	{bottomLeftHalf, bottomLeftTwoThirds, bottomLeftOneThirds},
	// 	{bottomRightHalf, bottomRightTwoThirds, bottomRightOneThirds}}
	// cornerFuncTurn := make([]int, len(cornerFuncs))

	cycleFuncs := func(funcs [][]resizeFunc, turns *[]int, i int) {
		hwnd := w32.GetForegroundWindow()
		if hwnd == 0 {
			panic("foreground window is NULL")
		}
		if lastResized != hwnd {
			*turns = make([]int, len(edgeFuncs)) // reset
		}
		if _, err := resize(hwnd, funcs[i][(*turns)[i]%len(funcs[i])]); err != nil {
			fmt.Printf("warn: resize: %v\n", err)
			return
		}
		(*turns)[i]++
		for j := 0; j < len(*turns); j++ {
			if j != i {
				(*turns)[j] = 0
			}
		}
	}

	// cycleEdgeFuncs := func(i int) { cycleFuncs(edgeFuncs, &edgeFuncTurn, i) }
	// cycleCornerFuncs := func(i int) { cycleFuncs(cornerFuncs, &cornerFuncTurn, i) }

	leftFuncs := [][]resizeFunc{{leftHalf, leftTwoThirds, leftOneThirds}}
	leftFuncTurn := make([]int, len(leftFuncs))
	cycleLeftFunc := func() { cycleFuncs(leftFuncs, &leftFuncTurn, 0) }

	rightFuncs := [][]resizeFunc{{rightHalf, rightTwoThirds, rightOneThirds}}
	rightFuncTurn := make([]int, len(rightFuncs))
	cycleRightFunc := func() { cycleFuncs(rightFuncs, &rightFuncTurn, 0) }

	topFuncs := [][]resizeFunc{{topMaximize, topHalf, bottomHalf}}
	topFuncTurn := make([]int, len(topFuncs))
	cycleTopFunc := func() { cycleFuncs(topFuncs, &topFuncTurn, 0) }

	hks := []HotKey{
		(HotKey{id: 1, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_LEFT, callback: func() { cycleLeftFunc() }}),
		(HotKey{id: 2, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_RIGHT, callback: func() { cycleRightFunc() }}),
		(HotKey{id: 3, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_UP, callback: func() { cycleTopFunc() }}),

		// (HotKey{id: 1, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_LEFT, callback: func() { cycleEdgeFuncs(0) }}),
		// (HotKey{id: 2, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_RIGHT, callback: func() { cycleEdgeFuncs(1) }}),
		// (HotKey{id: 3, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_UP, callback: func() { cycleEdgeFuncs(2) }}),
		// (HotKey{id: 4, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_DOWN, callback: func() { cycleEdgeFuncs(3) }}),
		// (HotKey{id: 5, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_LEFT, callback: func() { cycleCornerFuncs(0) }}),
		// (HotKey{id: 6, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_UP, callback: func() { cycleCornerFuncs(1) }}),
		// (HotKey{id: 7, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_DOWN, callback: func() { cycleCornerFuncs(2) }}),
		// (HotKey{id: 8, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_RIGHT, callback: func() { cycleCornerFuncs(3) }}),
		// (HotKey{id: 50, mod: MOD_SHIFT | MOD_WIN, vk: 0x46 /*F*/, callback: func() {
		// 	lastResized = 0 // cause edgeFuncTurn to be reset
		// 	if err := maximize(); err != nil {
		// 		fmt.Printf("warn: maximize: %v\n", err)
		// 		return
		// 	}
		// }}),
		// (HotKey{id: 60, mod: MOD_ALT | MOD_WIN, vk: 0x43 /*C*/, callback: func() {
		// 	lastResized = 0 // cause edgeFuncTurn to be reset
		// 	if _, err := resize(w32.GetForegroundWindow(), center); err != nil {
		// 		fmt.Printf("warn: resize: %v\n", err)
		// 		return
		// 	}
		// }}),
		// (HotKey{id: 70, mod: MOD_ALT | MOD_WIN, vk: 0x41 /*A*/, callback: func() {
		// 	hwnd := w32.GetForegroundWindow()
		// 	if err := toggleAlwaysOnTop(hwnd); err != nil {
		// 		fmt.Printf("warn: toggleAlwaysOnTop: %v\n", err)
		// 		return
		// 	}
		// 	fmt.Printf("> toggled always on top: %v\n", hwnd)
		// }}),
	}

	var failedHotKeys []HotKey
	for _, hk := range hks {
		if !RegisterHotKey(hk) {
			failedHotKeys = append(failedHotKeys, hk)
		}
	}
	if len(failedHotKeys) > 0 {
		msg := "The following hotkey(s) are in use by another process:\n\n"
		for _, hk := range failedHotKeys {
			msg += "  - " + hk.Describe() + "\n"
		}
		msg += "\nTo use these hotkeys in RectangleWin, close the other process using the key combination(s)."
		showMessageBox(msg)
	}

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, os.Interrupt)
	go func() {
		<-exitCh
		fmt.Println("exit signal received")
		systray.Quit() // causes WM_CLOSE, WM_QUIT, not sure if a side-effect
	}()

	// TODO systray/systray.go already locks the OS thread in init()
	// however it's not clear if GetMessage(0,0) will continue to work
	// as we run "go initTray()" and not pin the thread that initializes the
	// tray.
	initTray()
	if err := msgLoop(); err != nil {
		panic(err)
	}
}

func showMessageBox(text string) {
	w32.MessageBox(w32.GetActiveWindow(), text, "RectangleWin", w32.MB_ICONWARNING|w32.MB_OK)
}

type resizeFunc func(disp, cur w32.RECT) w32.RECT

func center(disp, cur w32.RECT) w32.RECT {
	// TODO find a way to round up divisions consistently as it causes multiple runs to shift by 1px
	w := (disp.Width() - cur.Width()) / 2
	h := (disp.Height() - cur.Height()) / 2
	return w32.RECT{
		Left:   disp.Left + w,
		Right:  disp.Left + w + cur.Width(),
		Top:    disp.Top + h,
		Bottom: disp.Top + h + cur.Height()}
}

func resize(hwnd w32.HWND, f resizeFunc) (bool, error) {
	if !isZonableWindow(hwnd) {
		fmt.Printf("warn: non-zonable window: %s\n", w32.GetWindowText(hwnd))
		return false, nil
	}
	rect := w32.GetWindowRect(hwnd)
	mon := w32.MonitorFromWindow(hwnd, w32.MONITOR_DEFAULTTONEAREST)
	hdc := w32.GetDC(hwnd)
	displayDPI := w32.GetDeviceCaps(hdc, w32.LOGPIXELSY)
	if !w32.ReleaseDC(hwnd, hdc) {
		return false, fmt.Errorf("failed to ReleaseDC:%d", w32.GetLastError())
	}
	var monInfo w32.MONITORINFO
	if !w32.GetMonitorInfo(mon, &monInfo) {
		return false, fmt.Errorf("failed to GetMonitorInfo:%d", w32.GetLastError())
	}

	ok, frame := w32.DwmGetWindowAttributeEXTENDED_FRAME_BOUNDS(hwnd)
	if !ok {
		return false, fmt.Errorf("failed to DwmGetWindowAttributeEXTENDED_FRAME_BOUNDS:%d", w32.GetLastError())
	}
	windowDPI := w32ex.GetDpiForWindow(hwnd)
	resizedFrame := resizeForDpi(frame, int32(windowDPI), int32(displayDPI))

	fmt.Printf("> window: 0x%x %#v (w:%d,h:%d) mon=0x%X(@ display DPI:%d)\n", hwnd, rect, rect.Width(), rect.Height(), mon, displayDPI)
	fmt.Printf("> DWM frame:        %#v (W:%d,H:%d) @ window DPI=%v\n", frame, frame.Width(), frame.Height(), windowDPI)
	fmt.Printf("> DPI-less frame:   %#v (W:%d,H:%d)\n", resizedFrame, resizedFrame.Width(), resizedFrame.Height())

	// calculate how many extra pixels go to win10 invisible borders
	lExtra := resizedFrame.Left - rect.Left
	rExtra := -resizedFrame.Right + rect.Right
	tExtra := resizedFrame.Top - rect.Top
	bExtra := -resizedFrame.Bottom + rect.Bottom

	newPos := f(monInfo.RcWork, resizedFrame)

	// adjust offsets based on invisible borders
	newPos.Left -= lExtra
	newPos.Top -= tExtra
	newPos.Right += rExtra
	newPos.Bottom += bExtra

	lastResized = hwnd
	if sameRect(rect, &newPos) {
		fmt.Println("no resize")
		return false, nil
	}

	fmt.Printf("> resizing to: %#v (W:%d,H:%d)\n", newPos, newPos.Width(), newPos.Height())
	if !w32.ShowWindow(hwnd, w32.SW_SHOWNORMAL) { // normalize window first if it's set to SW_SHOWMAXIMIZE (and therefore stays maximized)
		return false, fmt.Errorf("failed to normalize window ShowWindow:%d", w32.GetLastError())
	}
	if !w32.SetWindowPos(hwnd, 0, int(newPos.Left), int(newPos.Top), int(newPos.Width()), int(newPos.Height()), w32.SWP_NOZORDER|w32.SWP_NOACTIVATE) {
		return false, fmt.Errorf("failed to SetWindowPos:%d", w32.GetLastError())
	}
	rect = w32.GetWindowRect(hwnd)
	fmt.Printf("> post-resize: %#v(W:%d,H:%d)\n", rect, rect.Width(), rect.Height())
	return true, nil
}

func maximize() error {
	hwnd := w32.GetForegroundWindow()
	if !isZonableWindow(hwnd) {
		return errors.New("foreground window is not zonable")
	}
	if !w32.ShowWindow(hwnd, w32.SW_MAXIMIZE) {
		return fmt.Errorf("failed to ShowWindow:%d", w32.GetLastError())
	}
	return nil
}

func toggleAlwaysOnTop(hwnd w32.HWND) error {
	if !isZonableWindow(hwnd) {
		return errors.New("foreground window is not zonable")
	}

	if w32.GetWindowLong(hwnd, w32.GWL_EXSTYLE)&w32.WS_EX_TOPMOST != 0 {
		if !w32.SetWindowPos(hwnd, w32.HWND_NOTOPMOST, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE) {
			return fmt.Errorf("failed to SetWindowPos(HWND_NOTOPMOST): %v", w32.GetLastError())
		}
	} else {
		if !w32.SetWindowPos(hwnd, w32.HWND_TOPMOST, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE) {
			return fmt.Errorf("failed to SetWindowPos(HWND_TOPMOST) :%v", w32.GetLastError())
		}
	}
	return nil
}

func resizeForDpi(src w32.RECT, from, to int32) w32.RECT {
	return w32.RECT{
		Left:   src.Left * to / from,
		Right:  src.Right * to / from,
		Top:    src.Top * to / from,
		Bottom: src.Bottom * to / from,
	}
}

func sameRect(a, b *w32.RECT) bool {
	return a != nil && b != nil && reflect.DeepEqual(*a, *b)
}
