// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"math"

	"github.com/szabba/tob-cob/game"
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
)

type Grid struct {
	CellWidth  float64
	CellHeight float64
}

func (grid Grid) Matrix(col, row int) geometry.Mat {
	x, y := float64(col), float64(row)
	dx := grid.CellWidth * x
	dy := grid.CellHeight * y
	dr := geometry.V(dx, dy)
	return geometry.Translation(dr)
}

func (grid Grid) UnderCursor(src input.Source, cam Camera) game.Point {
	onScreen := src.MousePosition()
	toWorld := cam.Matrix(src.Bounds()).Invert()
	inWorld := toWorld.Apply(onScreen)
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
