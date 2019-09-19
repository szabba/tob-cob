// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package draw

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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

func (grid Grid) DrawCell(dst pixel.Target, col, row int, color pixel.RGBA) {
	imd := imdraw.New(nil)
	imd.Color = color
	imd.SetMatrix(grid.Matrix(col, row))
	wHalf := float64(grid.CellWidth) / 2
	imd.Push(pixel.V(-wHalf, 0), pixel.V(wHalf, grid.CellHeight))
	imd.Rectangle(1)
	imd.Draw(dst)
}
