// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"math"

	"github.com/faiface/pixel"

	"github.com/szabba/tob-cob/game"
)

type Grid struct {
	CellWidth  float64
	CellHeight float64
}

func (grid Grid) Matrix(col, row int) pixel.Matrix {
	x, y := float64(col), float64(row)
	dx := grid.CellWidth * x
	dy := grid.CellHeight * y
	dr := pixel.V(dx, dy)
	return pixel.IM.Moved(dr)
}

func (grid Grid) UnderCursor(input Input, cam Camera) game.Point {
	onScreen := input.MousePosition()
	inWorld := cam.Matrix(input.Bounds()).Unproject(onScreen)
	column := grid.underCursor(inWorld.X, grid.CellWidth)
	row := grid.underCursor(inWorld.Y, grid.CellWidth)
	return game.P(row, column)
}

func (grid Grid) underCursor(cursor, dim float64) int {
	if cursor < 0 {
		return -grid.underCursor(-cursor, dim)
	}

	cell := math.Floor((cursor + dim/2) / dim)
	return int(cell)
}
