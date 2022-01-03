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
	hotkeyRegistrations = make(map[int]*HotKey)
)

type HotKey struct {
	id, mod, vk int
	callback    func()
}

func (h HotKey) String() string { return fmt.Sprintf("mod=0x%x,vk=%d", h.mod, h.vk) }

func (h HotKey) Describe() string {
	var out string
	if h.mod&MOD_WIN == MOD_WIN {
		out += modKeyNames[MOD_WIN] + " + "
	}
	if h.mod&MOD_CONTROL == MOD_CONTROL {
		out += modKeyNames[MOD_CONTROL] + " + "
	}
	if h.mod&MOD_ALT == MOD_ALT {
		out += modKeyNames[MOD_ALT] + " + "
	}
	if h.mod&MOD_SHIFT == MOD_SHIFT {
		out += modKeyNames[MOD_SHIFT] + " + "
	}
	if v, ok := keyNames[h.vk]; ok {
		out += v
	} else {
		out += fmt.Sprintf("UNKNOWN KEY(0x%x)", h.vk)
	}
	return out
}

func RegisterHotKey(h HotKey) bool {
	if _, ok := hotkeyRegistrations[h.id]; ok {
		panic("hotkey id already registered") // TODO ok for now
	}
	ok := w32ex.RegisterHotKey(0, h.id, h.mod, h.vk)
	if ok {
		hotkeyRegistrations[h.id] = &h
	}
	return ok
}

func msgLoop() error {
	defer fmt.Println("event loop finished")
	for {
		var m w32.MSG
		c := w32.GetMessage(&m, 0, 0, 0)
		if c == -1 {
			return fmt.Errorf("GetMessage failed: %d", c)
		} else if c == 0 {
			// WM_QUIT received
			return nil
		}
		if m.Message == w32.WM_HOTKEY {
			h, ok := hotkeyRegistrations[int(m.WParam)]
			if !ok {
				return fmt.Errorf("hotkey without callback: %#v", m)
			}
			fmt.Printf("trace: hotkey id=%d (%s)\n", m.WParam, h)
			h.callback()
		} else {
			fmt.Printf("unhandled message received:0x%x %d\n", m.Message, m.Message)
			w32.TranslateMessage(&m)
			w32.DispatchMessage(&m)
		}
	}
}
