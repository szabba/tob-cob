// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"

	"github.com/szabba/tob-cob/game"
)

type GridOutline struct {
	Space   *game.Space
	Grid    Grid
	Color   pixel.RGBA
	Margins Margins
}

type Margins struct{ X, Y float64 }

func (o GridOutline) Draw(dst pixel.Target) {
	min, max := o.Space.Min(), o.Space.Max()
	var pt game.Point
	for pt.Row = min.Row; pt.Row <= max.Row; pt.Row++ {
		for pt.Column = min.Column; pt.Column <= max.Column; pt.Column++ {
			if !o.Space.At(pt).Exists() {
				continue
			}
			o.drawCell(dst, pt)
		}
	}
}

func (o GridOutline) drawCell(dst pixel.Target, pt game.Point) {
	matrix := o.cellMatrix(pt)

	halfWidth := o.Grid.CellWidth/2 - math.Abs(o.Margins.X)
	halfHeight := o.Grid.CellHeight/2 - math.Abs(o.Margins.Y)

	lowerLeft := pixel.V(-halfWidth, -halfHeight)
	upperRight := pixel.V(halfWidth, halfHeight)

	imd := imdraw.New(nil)
	imd.SetMatrix(matrix)
	imd.Color = o.Color
	imd.Push(lowerLeft, upperRight)
	imd.Rectangle(1)
	imd.Draw(dst)
}

func (o GridOutline) cellMatrix(pt game.Point) pixel.Matrix {
	return o.Grid.Matrix(pt.Column, pt.Row)
}
