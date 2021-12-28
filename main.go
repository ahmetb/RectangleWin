package main

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gonutz/w32/v2"
	"golang.org/x/sys/windows"
)

var lastResized w32.HWND

func main() {
	w32.EnumWindows(func(window w32.HWND) bool {
		fmt.Println("->", window, w32.GetWindowText(window))
		fmt.Printf("\t%#v\n", w32.GetWindowRect(window))
		return true

	})

	EnumMonitors(func(d w32.HMONITOR) bool {
		var v w32.MONITORINFO
		if !w32.GetMonitorInfo(d, &v) {
			return false
		}
		fmt.Printf("> monitor:0x%x\n", d)
		fmt.Printf("     rcwork:%#v\n", v.RcWork)
		fmt.Printf("  rcmonitor:%#v\n", v.RcMonitor)
		fmt.Printf("    primary:%#v\n", v.DwFlags&w32.MONITORINFOF_PRIMARY > 0)

		ok, n := w32.GetNumberOfPhysicalMonitorsFromHMONITOR(d)
		if !ok {
			fmt.Printf("  physical monitors: failed to query count: %d\n", w32.GetLastError())
		} else {
			fmt.Printf("  physical monitors: %d\n", n)
			pMon := make([]w32.PHYSICAL_MONITOR, n)
			if !w32.GetPhysicalMonitorsFromHMONITOR(d, pMon) {
				fmt.Printf("  physical monitors: failed to get physical monitors: %d\n", w32.GetLastError())
			} else {
				for i, p := range pMon {
					name := windows.UTF16ToString(p.Description[:])
					fmt.Printf("  physical monitor#%d: %s\n", i, name)
				}
			}
		}
		return true
	})

	resizeFuncs := [][]resizeFunc{
		{leftHalf, leftTwoThirds, leftOneThirds},
		{rightHalf, rightTwoThirds, rightOneThirds},
		{topHalf, topTwoThirds, topOneThirds},
		{bottomHalf, bottomTwoThirds, bottomOneThirds}}
	curFuncs := make([]int, len(resizeFuncs))
	applyResize := func(i int) {
		hand := w32.GetForegroundWindow()
		if hand == 0 {
			panic("foreground window is NULL")
		}
		if lastResized != hand {
			curFuncs = make([]int, len(resizeFuncs)) // reset
		}
		if _, err := resize(hand, resizeFuncs[i][curFuncs[i]%len(resizeFuncs[i])]); err != nil {
			panic(err)
		}
		curFuncs[i]++
		for j := 0; j < len(curFuncs); j++ {
			if j != i {
				curFuncs[j] = 0
			}
		}
	}

	RegisterHotKey(HotKey{id: 1, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_LEFT, callback: func() { applyResize(0) }})
	RegisterHotKey(HotKey{id: 2, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_RIGHT, callback: func() { applyResize(1) }})
	RegisterHotKey(HotKey{id: 3, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_UP, callback: func() { applyResize(2) }})
	RegisterHotKey(HotKey{id: 4, mod: MOD_ALT | MOD_WIN | MOD_NOREPEAT, vk: w32.VK_DOWN, callback: func() { applyResize(3) }})
	RegisterHotKey(HotKey{id: 5, mod: MOD_SHIFT | MOD_WIN, vk: 0x46 /*F*/, callback: func() {
		lastResized = 0 // cause curFuncs to be reset
		if err := maximize(); err != nil {
			panic(err)
		}
	}})
	// TODO center
	if err := StartHotKeyListen(); err != nil {
		// TODO reset curFuncs
		panic(err)
	}
}

type resizeFunc func(display w32.RECT) w32.RECT

func toLeft(d w32.RECT, mul, div int32) w32.RECT {
	return w32.RECT{Left: 0, Top: 0, Right: (d.Width() * mul) / div, Bottom: d.Height()}
}

func toRight(d w32.RECT, mul, div int32) w32.RECT {
	return w32.RECT{Left: d.Width() - d.Width()*mul/div, Top: 0, Right: d.Width(), Bottom: d.Height()}
}

func toTop(d w32.RECT, mul, div int32) w32.RECT {
	return w32.RECT{Left: 0, Top: 0, Right: d.Width(), Bottom: d.Height() * mul / div}
}

func toBottom(d w32.RECT, mul, div int32) w32.RECT {
	return w32.RECT{Left: 0, Top: d.Height() - d.Height()*mul/div, Right: d.Width(), Bottom: d.Height()}
}

func leftHalf(disp w32.RECT) w32.RECT      { return toLeft(disp, 1, 2) }
func leftOneThirds(disp w32.RECT) w32.RECT { return toLeft(disp, 1, 3) }
func leftTwoThirds(disp w32.RECT) w32.RECT { return toLeft(disp, 2, 3) }

func topHalf(disp w32.RECT) w32.RECT      { return toTop(disp, 1, 2) }
func topOneThirds(disp w32.RECT) w32.RECT { return toTop(disp, 1, 3) }
func topTwoThirds(disp w32.RECT) w32.RECT { return toTop(disp, 2, 3) }

func rightHalf(disp w32.RECT) w32.RECT      { return toRight(disp, 1, 2) }
func rightOneThirds(disp w32.RECT) w32.RECT { return toRight(disp, 1, 3) }
func rightTwoThirds(disp w32.RECT) w32.RECT { return toRight(disp, 2, 3) }

func bottomHalf(disp w32.RECT) w32.RECT      { return toBottom(disp, 1, 2) }
func bottomOneThirds(disp w32.RECT) w32.RECT { return toBottom(disp, 1, 3) }
func bottomTwoThirds(disp w32.RECT) w32.RECT { return toBottom(disp, 2, 3) }

func resize(hand w32.HWND, f resizeFunc) (bool, error) {
	rect := w32.GetWindowRect(hand)
	title := w32.GetWindowText(hand)
	if title == "Program Manager" {
		return false, errors.New("foreground window is ProgMan")
	}

	mon := w32.MonitorFromWindow(hand, w32.MONITOR_DEFAULTTONULL)
	hdc := w32.GetDC(hand)
	monDPI := w32.GetDeviceCaps(hdc, w32.LOGPIXELSY)
	if !w32.ReleaseDC(hand, hdc) {
		return false, fmt.Errorf("failed to ReleaseDC:%d", w32.GetLastError())
	}
	windowDPI := GetDpiForWindow(hand)
	var monInfo w32.MONITORINFO
	w32.GetMonitorInfo(mon, &monInfo)

	ok, frame := w32.DwmGetWindowAttributeEXTENDED_FRAME_BOUNDS(hand)
	if !ok {
		return false, fmt.Errorf("failed to DwmGetWindowAttributeEXTENDED_FRAME_BOUNDS:%d", w32.GetLastError())
	}
	resizedFrame := resizeForDpi(frame, int32(windowDPI), int32(monDPI))

	fmt.Printf("> window:        %#v (w:%d,h:%d) mon=0x%X(@ DPI:%d)\n", rect, rect.Width(), rect.Height(), mon, monDPI)
	fmt.Printf("> DWM frame:     %#v (W:%d,H:%d) @ DPI=%v\n", frame, frame.Width(), frame.Height(), windowDPI)
	fmt.Printf("> DPI-less frame: %#v (W:%d,H:%d)\n", resizedFrame, resizedFrame.Width(), resizedFrame.Height())

	// calculate how many extra pixels go to win10 invisible borders
	lExtra := resizedFrame.Left - rect.Left
	rExtra := -resizedFrame.Right + rect.Right
	tExtra := resizedFrame.Top - rect.Top
	bExtra := -resizedFrame.Bottom + rect.Bottom

	newPos := f(monInfo.RcWork)

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

func GetDpiForWindow(hwnd w32.HWND) int32 {
	r1, _, _ := user32.NewProc("GetDpiForWindow").Call(uintptr(hwnd))
	return int32(r1)
}

func sameRect(a, b *w32.RECT) bool {
	return a != nil && b != nil && reflect.DeepEqual(*a, *b)
}
