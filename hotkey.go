package main

import (
	"fmt"

	"github.com/ahmetb/RectangleWin/w32ex"
	"github.com/gonutz/w32/v2"
)

var (
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

func (h HotKey) String() string { return fmt.Sprintf("mod=0x%x,vk=%d", h.mod, h.vk) }

func RegisterHotKey(h HotKey) {
	// TODO not safe for concurrent modification
	if _, ok := hotkeys[h.id]; ok {
		panic("hotkey id already registered") // TODO ok for now
	}
	hotkeys[h.id] = &h
	w32ex.RegisterHotKey(0, h.id, h.mod, h.vk)
}

func hotKeyLoop() error {
	for {
		var m w32.MSG
		if c := w32.GetMessage(&m, 0, w32.WM_HOTKEY, w32.WM_HOTKEY); c <= 0 {
			return fmt.Errorf("GetMessage failed: %d", c)
		}
		h, ok := hotkeys[int(m.WParam)]
		if !ok {
			return fmt.Errorf("hotkey without callback: %#v", m)
		}
		fmt.Printf("trace: hotkey id=%d (%s)\n", m.WParam, h)
		h.callback()
	}
}
