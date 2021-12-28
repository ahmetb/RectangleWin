package main

import "github.com/gonutz/w32/v2"

// TODO find a way to round up divisions consistently, otherwise we end up with off by 1px

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

func topLeftHalf(disp, _ w32.RECT) w32.RECT      { return merge(toLeft(disp, 1, 2), toTop(disp, 1, 2)) }
func topLeftTwoThirds(disp, _ w32.RECT) w32.RECT { return merge(toLeft(disp, 2, 3), toTop(disp, 1, 2)) }
func topLeftOneThirds(disp, _ w32.RECT) w32.RECT { return merge(toLeft(disp, 1, 3), toTop(disp, 1, 2)) }

func topRightHalf(disp, _ w32.RECT) w32.RECT { return merge(toRight(disp, 1, 2), toTop(disp, 1, 2)) }
func topRightTwoThirds(disp, _ w32.RECT) w32.RECT {
	return merge(toRight(disp, 2, 3), toTop(disp, 1, 2))
}
func topRightOneThirds(disp, _ w32.RECT) w32.RECT {
	return merge(toRight(disp, 1, 3), toTop(disp, 1, 2))
}

func bottomLeftHalf(disp, _ w32.RECT) w32.RECT {
	return merge(toLeft(disp, 1, 2), toBottom(disp, 1, 2))
}
func bottomLeftTwoThirds(disp, _ w32.RECT) w32.RECT {
	return merge(toLeft(disp, 2, 3), toBottom(disp, 1, 2))
}
func bottomLeftOneThirds(disp, _ w32.RECT) w32.RECT {
	return merge(toLeft(disp, 1, 3), toBottom(disp, 1, 2))
}

func bottomRightHalf(disp, _ w32.RECT) w32.RECT {
	return merge(toRight(disp, 1, 2), toBottom(disp, 1, 2))
}
func bottomRightTwoThirds(disp, _ w32.RECT) w32.RECT {
	return merge(toRight(disp, 2, 3), toBottom(disp, 1, 2))
}
func bottomRightOneThirds(disp, _ w32.RECT) w32.RECT {
	return merge(toRight(disp, 1, 3), toBottom(disp, 1, 2))
}

func merge(x, y w32.RECT) w32.RECT {
	return w32.RECT{Left: x.Left, Right: x.Right, Top: y.Top, Bottom: y.Bottom}
}
