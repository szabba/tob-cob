// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"math"

	"github.com/szabba/tob-cob/game/grid"
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
)

type GridDimensions struct {
	CellWidth  float64
	CellHeight float64
}

func (d GridDimensions) Matrix(col, row int) geometry.Mat {
	x, y := float64(col), float64(row)
	dx := d.CellWidth * x
	dy := d.CellHeight * y
	dr := geometry.V(dx, dy)
	return geometry.Translation(dr)
}

func (d GridDimensions) UnderCursor(src input.Source, cam Camera) grid.Point {
	onScreen := src.MousePosition()
	toWorld := cam.Matrix(src.Bounds()).Invert()
	inWorld := toWorld.Apply(onScreen)
	column := d.underCursor(inWorld.X, d.CellWidth)
	row := d.underCursor(inWorld.Y, d.CellWidth)
	return grid.P(row, column)
}

func (d GridDimensions) underCursor(cursor, dim float64) int {
	if cursor < 0 {
		return -d.underCursor(-cursor, dim)
	}

	cell := math.Floor((cursor + dim/2) / dim)
	return int(cell)
}
