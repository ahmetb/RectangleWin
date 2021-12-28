package main

import (
	"syscall"

	"github.com/gonutz/w32/v2"
)

func EnumMonitors(f func(d w32.HMONITOR) bool) bool {
	callback := syscall.NewCallback(func(h, _, _, _ uintptr) uintptr {
		if f(w32.HMONITOR(h)) {
			return 1
		}
		return 0
	})
	return w32.EnumDisplayMonitors(0, nil, callback, 0)
}
