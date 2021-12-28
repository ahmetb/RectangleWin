package main

import (
	"fmt"
	"syscall"

	"github.com/gonutz/w32/v2"
)

var (
	user32  = syscall.NewLazyDLL("user32.dll")
	hotkeys = make(map[int]*HotKey)
)

const (
	MOD_ALT      = 0x0001
	MOD_CONTROL  = 0x0002
	MOD_NOREPEAT = 0x4000
	MOD_SHIFT    = 0x0004
	MOD_WIN      = 0x0008
)

type HotKey struct {
	id, mod, vk int
	callback    func()
}

func registerHotKey(hwnd w32.HWND, id, mod, vk int) {
	r1, _, _ := user32.NewProc("RegisterHotKey").Call(uintptr(hwnd), uintptr(id), uintptr(mod), uintptr(vk))
	if r1 == 0 {
		panic(fmt.Errorf("failed to register hotkey:%d lastErr:%d", r1, w32.GetLastError()))
	}
}

func RegisterHotKey(h HotKey) {
	// TODO not safe for concurrent modification
	if _, ok := hotkeys[h.id]; ok {
		panic("hotkey id already registered") // TODO ok for now
	}
	hotkeys[h.id] = &h
	registerHotKey(0, h.id, h.mod, h.vk)
}

func StartHotKeyListen() error {
	for {
		var m w32.MSG
		if c := w32.GetMessage(&m, 0, w32.WM_HOTKEY, w32.WM_HOTKEY); c <= 0 {
			return fmt.Errorf("GetMessage failed: %d", c)
		}
		h, ok := hotkeys[int(m.WParam)]
		if !ok {
			return fmt.Errorf("hotkey without callback: %#v", m)
		}
		fmt.Printf("trace: hotkey %d (mod=0x%X,vk=%d)\n", m.WParam, h.mod, h.vk)
		h.callback()
	}
}
