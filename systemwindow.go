// Copyright 2022 Ahmet Alp Balkan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"strings"

	"github.com/ahmetb/RectangleWin/w32ex"
	"github.com/gonutz/w32/v2"
)

const (
	GWL_EXSTYLE = -20
	GWL_STYLE   = -16
)

func isZonableWindow(hwnd w32.HWND) bool {
	if hwnd == 0 {
		return false
	}
	return isStandardWindow(hwnd) && hasNoVisibleOwner(hwnd)
}

func hasNoVisibleOwner(hwnd w32.HWND) bool {
	owner := w32.GetWindow(hwnd, w32.GW_OWNER)
	if owner == 0 {
		return true
	}
	if !w32.IsWindowVisible(owner) {
		return true
	}
	rect := w32.GetWindowRect(owner)
	if rect == nil {
		return false
	}
	return rect.Width() == 0 || rect.Height() == 0
}

func isStandardWindow(hwnd w32.HWND) bool {
	// adapted from https://github.com/microsoft/PowerToys/blob/7d0304fd06939d9f552e75be9c830db22f8ff9e2/src/modules/fancyzones/FancyZonesLib/util.cpp#L403
	if w32ex.GetAncestor(hwnd, w32ex.GA_ROOT) != hwnd ||
		!w32.IsWindowVisible(hwnd) {
		return false
	}

	for _, sysWindow := range []w32.HWND{w32.GetDesktopWindow(), w32ex.GetShellWindow()} {
		if hwnd == sysWindow {
			return false
		}
	}

	style := w32.GetWindowLong(hwnd, GWL_STYLE)
	// a window with think frame and minimize/maximize buttons
	if uint32(style)&w32.WS_POPUP == w32.WS_POPUP &&
		style&w32.WS_THICKFRAME == w32.WS_THICKFRAME &&
		style&w32.WS_MINIMIZEBOX == 0 &&
		style&w32.WS_MAXIMIZEBOX == 0 {
		return false
	}
	exStyle := w32.GetWindowLong(hwnd, GWL_EXSTYLE)
	if uint32(style)&w32.WS_CHILD == w32.WS_CHILD ||
		style&w32.WS_DISABLED == w32.WS_DISABLED ||
		exStyle&w32.WS_EX_TOOLWINDOW == w32.WS_EX_TOOLWINDOW ||
		exStyle&w32.WS_EX_NOACTIVATE == w32.WS_EX_NOACTIVATE {
		return false
	}

	className, ok := w32.GetClassName(hwnd)
	if !ok {
		panic("GetClassName failed")
	}
	return !isSystemClassName(className)
}

func isSystemClassName(className string) bool {
	// adapted from https://github.com/microsoft/PowerToys/blob/7d0304fd06939d9f552e75be9c830db22f8ff9e2/tools/FancyZones_zonable_tester/main.cpp#L135
	for _, c := range []string{
		"SysListView32",
		"WorkerW",
		"Shell_TrayWnd",
		"Shell_SecondaryTrayWnd",
		"Progman",
	} {
		if strings.EqualFold(c, className) {
			return true
		}
	}
	return false
}
