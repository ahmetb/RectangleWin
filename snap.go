package main

import "github.com/gonutz/w32/v2"

func toLeft(d w32.RECT, mul, div int32) w32.RECT {
	// TODO find a way to round up divisions consistently
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

func leftHalf(disp, _ w32.RECT) w32.RECT      { return toLeft(disp, 1, 2) }
func leftOneThirds(disp, _ w32.RECT) w32.RECT { return toLeft(disp, 1, 3) }
func leftTwoThirds(disp, _ w32.RECT) w32.RECT { return toLeft(disp, 2, 3) }

func topHalf(disp, _ w32.RECT) w32.RECT      { return toTop(disp, 1, 2) }
func topOneThirds(disp, _ w32.RECT) w32.RECT { return toTop(disp, 1, 3) }
func topTwoThirds(disp, _ w32.RECT) w32.RECT { return toTop(disp, 2, 3) }

func rightHalf(disp, _ w32.RECT) w32.RECT      { return toRight(disp, 1, 2) }
func rightOneThirds(disp, _ w32.RECT) w32.RECT { return toRight(disp, 1, 3) }
func rightTwoThirds(disp, _ w32.RECT) w32.RECT { return toRight(disp, 2, 3) }

func bottomHalf(disp, _ w32.RECT) w32.RECT      { return toBottom(disp, 1, 2) }
func bottomOneThirds(disp, _ w32.RECT) w32.RECT { return toBottom(disp, 1, 3) }
func bottomTwoThirds(disp, _ w32.RECT) w32.RECT { return toBottom(disp, 2, 3) }
