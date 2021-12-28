package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gonutz/w32/v2"

	"github.com/ahmetb/RectangleWin/w32ex"
)

var lastResized w32.HWND

func main() {
	edgeFuncs := [][]resizeFunc{
		{leftHalf, leftTwoThirds, leftOneThirds},
		{rightHalf, rightTwoThirds, rightOneThirds},
		{topHalf, topTwoThirds, topOneThirds},
		{bottomHalf, bottomTwoThirds, bottomOneThirds}}
	edgeFuncTurn := make([]int, len(edgeFuncs))
	cornerFuncs := [][]resizeFunc{
		{topLeftHalf, topLeftTwoThirds, topLeftOneThirds},
		{topRightHalf, topRightTwoThirds, topRightOneThirds},
		{bottomLeftHalf, bottomLeftTwoThirds, bottomLeftOneThirds},
		{bottomRightHalf, bottomRightTwoThirds, bottomRightOneThirds}}
	cornerFuncTurn := make([]int, len(cornerFuncs))

	cycleFuncs := func(funcs [][]resizeFunc, turns *[]int, i int) {
		hand := w32.GetForegroundWindow()
		if hand == 0 {
			panic("foreground window is NULL")
		}
		if lastResized != hand {
			*turns = make([]int, len(edgeFuncs)) // reset
		}
		if _, err := resize(hand, funcs[i][(*turns)[i]%len(funcs[i])]); err != nil {
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

	cycleEdgeFuncs := func(i int) { cycleFuncs(edgeFuncs, &edgeFuncTurn, i) }
	cycleCornerFuncs := func(i int) { cycleFuncs(cornerFuncs, &cornerFuncTurn, i) }

	RegisterHotKey(HotKey{id: 1, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_LEFT, callback: func() { cycleEdgeFuncs(0) }})
	RegisterHotKey(HotKey{id: 2, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_RIGHT, callback: func() { cycleEdgeFuncs(1) }})
	RegisterHotKey(HotKey{id: 3, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_UP, callback: func() { cycleEdgeFuncs(2) }})
	RegisterHotKey(HotKey{id: 4, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_DOWN, callback: func() { cycleEdgeFuncs(3) }})

	RegisterHotKey(HotKey{id: 5, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_LEFT, callback: func() { cycleCornerFuncs(0) }})
	RegisterHotKey(HotKey{id: 6, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_UP, callback: func() { cycleCornerFuncs(1) }})
	RegisterHotKey(HotKey{id: 7, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_DOWN, callback: func() { cycleCornerFuncs(2) }})
	RegisterHotKey(HotKey{id: 8, mod: MOD_CONTROL | MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_RIGHT, callback: func() { cycleCornerFuncs(3) }})

	RegisterHotKey(HotKey{id: 50, mod: MOD_SHIFT | MOD_WIN, vk: 0x46 /*F*/, callback: func() {
		lastResized = 0 // cause edgeFuncTurn to be reset
		if err := maximize(); err != nil {
			panic(err)
		}
	}})
	RegisterHotKey(HotKey{id: 60, mod: MOD_ALT | MOD_WIN, vk: 0x43 /*C*/, callback: func() {
		lastResized = 0 // cause edgeFuncTurn to be reset
		// TODO find a common way to GetForegroundWindow and validate it
		if _, err := resize(w32.GetForegroundWindow(), center); err != nil {
			fmt.Printf("warn: resize: %v\n", err)
			return
		}
	}})
	if err := startHotKeyListen(); err != nil {
		panic(err)
	}
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

func resize(hand w32.HWND, f resizeFunc) (bool, error) {
	if isSystemWindow(hand) {
		return false, nil
	}
	rect := w32.GetWindowRect(hand)
	mon := w32.MonitorFromWindow(hand, w32.MONITOR_DEFAULTTONULL)
	hdc := w32.GetDC(hand)
	displayDPI := w32.GetDeviceCaps(hdc, w32.LOGPIXELSY)
	if !w32.ReleaseDC(hand, hdc) {
		return false, fmt.Errorf("failed to ReleaseDC:%d", w32.GetLastError())
	}
	var monInfo w32.MONITORINFO
	if !w32.GetMonitorInfo(mon, &monInfo) {
		return false, fmt.Errorf("failed to GetMonitorInfo:%d", w32.GetLastError())
	}

	ok, frame := w32.DwmGetWindowAttributeEXTENDED_FRAME_BOUNDS(hand)
	if !ok {
		return false, fmt.Errorf("failed to DwmGetWindowAttributeEXTENDED_FRAME_BOUNDS:%d", w32.GetLastError())
	}
	windowDPI := w32ex.GetDpiForWindow(hand)
	resizedFrame := resizeForDpi(frame, int32(windowDPI), int32(displayDPI))

	fmt.Printf("> window: 0x%x       %#v (w:%d,h:%d) mon=0x%X(@ DPI:%d)\n", hand, rect, rect.Width(), rect.Height(), mon, displayDPI)
	fmt.Printf("> DWM frame:     %#v (W:%d,H:%d) @ DPI=%v\n", frame, frame.Width(), frame.Height(), windowDPI)
	fmt.Printf("> DPI-less frame: %#v (W:%d,H:%d)\n", resizedFrame, resizedFrame.Width(), resizedFrame.Height())

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

	lastResized = hand
	if sameRect(rect, &newPos) {
		fmt.Println("no resize")
		return false, nil
	}

	fmt.Printf("> resizing to: %#v (W:%d,H:%d)\n", newPos, newPos.Width(), newPos.Height())
	if !w32.ShowWindow(hand, w32.SW_SHOWNORMAL) { // normalize window first if it's set to SW_SHOWMAXIMIZE (and therefore stays maximized)
		return false, fmt.Errorf("failed to normalize window ShowWindow:%d", w32.GetLastError())
	}
	if !w32.SetWindowPos(hand, 0, int(newPos.Left), int(newPos.Top), int(newPos.Width()), int(newPos.Height()), w32.SWP_NOZORDER|w32.SWP_NOACTIVATE) {
		return false, fmt.Errorf("failed to SetWindowPos:%d", w32.GetLastError())
	}
	rect = w32.GetWindowRect(hand)
	fmt.Printf("> post-resize: %#v(W:%d,H:%d)\n", rect, rect.Width(), rect.Height())
	return true, nil
}

func maximize() error {
	// TODO find a common way to GetForegroundWindow and validate it
	hwnd := w32.GetForegroundWindow()
	if hwnd == 0 {
		return errors.New("foreground window is NULL")
	}
	if !w32.ShowWindow(hwnd, w32.SW_MAXIMIZE) {
		return fmt.Errorf("failed to ShowWindow:%d", w32.GetLastError())
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

func isSystemWindow(hwnd w32.HWND) bool {
	// FIXME: this doesn't work for cmd, powershell or Windows Terminal app
	proc := w32ex.GetWindowModuleFileName(hwnd)
	return proc == ""
}
