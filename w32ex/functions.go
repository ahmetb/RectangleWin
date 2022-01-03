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

package w32ex

import (
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

const (
	GA_PARENT    = 1
	GA_ROOT      = 2
	GA_ROOTOWNER = 3
)

var user32 = syscall.NewLazyDLL("user32.dll")

func RegisterHotKey(hwnd w32.HWND, id, mod, vk int) bool {
	r1, _, _ := user32.NewProc("RegisterHotKey").Call(uintptr(hwnd), uintptr(id), uintptr(mod), uintptr(vk))
	return r1 != 0
}

func GetDpiForWindow(hwnd w32.HWND) int32 {
	r1, _, _ := user32.NewProc("GetDpiForWindow").Call(uintptr(hwnd))
	return int32(r1)
}

func GetWindowModuleFileName(hwnd w32.HWND) string {
	var path [32768]uint16
	ret, _, _ := user32.NewProc("GetWindowModuleFileNameW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&path[0])),
		uintptr(len(path)),
	)
	if ret == 0 {
		return ""
	}
	return syscall.UTF16ToString(path[:])
}

func GetAncestor(hwnd w32.HWND, gaFlags uint) w32.HWND {
	r1, _, _ := user32.NewProc("GetAncestor").Call(uintptr(hwnd), uintptr(gaFlags))
	return w32.HWND(r1)
}

func GetShellWindow() (hwnd w32.HWND) {
	r1, _, _ := user32.NewProc("GetShellWindow").Call()
	return w32.HWND(r1)
}

func SetProcessDPIAware() bool {
	r1, _, _ := user32.NewProc("SetProcessDPIAware").Call()
	return r1 != 0
}
