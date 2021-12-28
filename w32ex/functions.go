package w32ex

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

var user32 = syscall.NewLazyDLL("user32.dll")

func RegisterHotKey(hwnd w32.HWND, id, mod, vk int) {
	r1, _, _ := user32.NewProc("RegisterHotKey").Call(uintptr(hwnd), uintptr(id), uintptr(mod), uintptr(vk))
	if r1 == 0 {
		panic(fmt.Errorf("failed to register hotkey mod=0x%x,vk=%d err:%d lastErr:%d", mod, vk, r1, w32.GetLastError()))
	}
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
