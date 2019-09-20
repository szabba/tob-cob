// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Cell struct {
	Col, Row int
	Grid     Grid
}

func (cell Cell) Matrix() pixel.Matrix {
	return cell.Grid.Matrix(cell.Col, cell.Row)
}

func (cell Cell) Draw(dst pixel.Target) {
	imd := imdraw.New(nil)
	imd.SetMatrix(cell.Matrix())
	wHalf := float64(cell.Grid.CellWidth) / 2
	imd.Push(pixel.V(-wHalf, 0), pixel.V(wHalf, cell.Grid.CellHeight))

	imd.Color = pixel.RGB(1, 1, 1)
	imd.Rectangle(1)

	imd.Draw(dst)
}
