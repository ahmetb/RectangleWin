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

func hotKeyLoop(tray *NOTIFYICONDATA) error {
	for {
		var m w32.MSG
		tray.HWnd = uintptr(w32.GetForegroundWindow())
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
