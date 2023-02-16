// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pixelglinput

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
)

type Source struct {
	win *pixelgl.Window
}

var _ input.Source = Source{}

func New(w *pixelgl.Window) Source { return Source{w} }

func (src Source) Focused() bool {
	if src.win == nil {
		return false
	}

	return src.win.Focused()
}

func (src Source) Pressed(btn input.Button) bool {
	if src.win == nil {
		return false
	}

	glBtn, ok := btnMapping[btn]
	if !ok {
		return false
	}

	return src.win.Pressed(glBtn)
}

var btnMapping = map[input.Button]pixelgl.Button{
	input.KeyF(): pixelgl.KeyF,

	input.KeyLeft():  pixelgl.KeyLeft,
	input.KeyUp():    pixelgl.KeyUp,
	input.KeyRight(): pixelgl.KeyRight,
	input.KeyDown():  pixelgl.KeyDown,
}

func (src Source) Bounds() geometry.Rect {
	if src.win == nil {
		return geometry.Rect{}
	}

	bds := src.win.Bounds()
	return geometry.Rect{
		Min: geometry.Vec{bds.Min.X, bds.Min.Y},
		Max: geometry.Vec{bds.Max.X, bds.Max.Y},
	}
}

func (src Source) MouseInsideWindow() bool {
	if src.win == nil {
		return false
	}
	return src.win.MouseInsideWindow()
}

func (src Source) MousePosition() geometry.Vec {
	if src.win == nil {
		return geometry.Vec{}
	}
	return v2geo(src.win.MousePosition())
}

func v2geo(v pixel.Vec) geometry.Vec {
	return geometry.Vec{
		X: v.X,
		Y: v.Y,
	}
}
