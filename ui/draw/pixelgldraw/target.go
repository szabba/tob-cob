// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pixelgldraw

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/geometry"
)

type Target struct {
	win *pixelgl.Window
}

func New(win *pixelgl.Window) Target {
	return Target{win}
}

var _ draw.Target = Target{}

func (t Target) SetMatrix(m geometry.Mat) {
	matrix := toPxM(m)
	t.win.SetMatrix(matrix)
}

func (t Target) Rectangle(r geometry.Rect, m geometry.Mat, fill color.Color) {

	imd := imdraw.New(nil)
	imd.Color = fill

	matrix := toPxM(m)
	imd.SetMatrix(matrix)

	lowerLeft := pixel.V(r.Min.X, r.Min.Y)
	upperRight := pixel.V(r.Max.X, r.Max.Y)

	imd.Push(lowerLeft, upperRight)

	imd.Rectangle(1)
	imd.Draw(t.win)
}

func toPxM(m geometry.Mat) pixel.Matrix {
	return pixel.Matrix{m[0][0], m[1][0], m[0][1], m[1][1], m[0][2], m[1][2]}
}