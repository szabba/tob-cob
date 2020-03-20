// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"github.com/faiface/pixel"

	"github.com/szabba/tob-cob/game"
)

type Grid struct {
	CellWidth  float64
	CellHeight float64
	Dx, Dy     float64
}

func (grid Grid) Matrix(col, row int) pixel.Matrix {
	x, y := float64(col), float64(row)
	dx := grid.CellWidth*x + grid.Dx*x
	dy := grid.CellHeight*y + grid.Dy*y
	dr := pixel.V(dx, dy)
	return pixel.IM.Moved(dr)
}

func (grid Grid) Cell(col, row int) Cell {
	return Cell{col, row, grid}
}

type GridOutline struct {
	Space *game.Space
	Grid  Grid
	Color pixel.RGBA
}

func (o GridOutline) Draw(dst pixel.Target) {
	min, max := o.Space.Min(), o.Space.Max()
	var pt game.Point
	for pt.Row = min.Row; pt.Row <= max.Row; pt.Row++ {
		for pt.Column = min.Column; pt.Column <= max.Column; pt.Column++ {
			if !o.Space.At(pt).Exists() {
				continue
			}
			o.Grid.Cell(pt.Column, pt.Row).Draw(dst, o.Color)
		}
	}
}
