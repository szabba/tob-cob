// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui_test

import (
	"testing"

	"github.com/faiface/pixel"
	"github.com/szabba/assert"
	"github.com/szabba/tob-cob/game"
	"github.com/szabba/tob-cob/ui"
)

func TestGridMatrix(t *testing.T) {
	grid := ui.Grid{CellWidth: 20, CellHeight: 10}
	tests := map[string]struct {
		Column, Row int

		RelativeToCell  pixel.Vec
		RelativeToWorld pixel.Vec
	}{
		"MiddleOfOriginCell": {},

		"TopEdgeOfOriginCell": {
			RelativeToCell:  pixel.V(0, grid.CellHeight/2),
			RelativeToWorld: pixel.V(0, grid.CellHeight/2),
		},
		"BottomEdgeOfOriginCell": {
			RelativeToCell:  pixel.V(0, -grid.CellHeight/2),
			RelativeToWorld: pixel.V(0, -grid.CellHeight/2),
		},
		"LeftEdgeOfOriginCell": {
			RelativeToCell:  pixel.V(-grid.CellWidth/2, 0),
			RelativeToWorld: pixel.V(-grid.CellWidth/2, 0),
		},
		"RightEdgeOfOriginCell": {
			RelativeToCell:  pixel.V(grid.CellWidth/2, 0),
			RelativeToWorld: pixel.V(grid.CellWidth/2, 0),
		},

		"MiddleOfCellAboveOrigin": {
			Row:             1,
			RelativeToWorld: pixel.V(0, grid.CellHeight),
		},
		"MiddleOfCellBelowOrigin": {
			Row:             -1,
			RelativeToWorld: pixel.V(0, -grid.CellHeight),
		},
		"MiddleOfCellRightOfOrigin": {
			Column:          1,
			RelativeToWorld: pixel.V(grid.CellWidth, 0),
		},
		"MiddleOfCellLeftOfOrigin": {
			Column:          -1,
			RelativeToWorld: pixel.V(-grid.CellWidth, 0),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			matrix := grid.Matrix(tt.Column, tt.Row)

			// then
			relativeToWorld := matrix.Project(tt.RelativeToCell)
			assert.That(
				relativeToWorld == tt.RelativeToWorld,
				t.Errorf, "got %#v world-coordinate, want %#v", relativeToWorld, tt.RelativeToWorld)
		})
	}
}

func TestGridUnderCursor(t *testing.T) {
	grid := ui.Grid{CellWidth: 20, CellHeight: 10}
	tests := map[string]struct {
		LookingAt pixel.Vec
		MouseAt   func() pixel.Vec

		Cell game.Point
	}{
		"AtOrigin": {},

		"AtLeftEdgeOfOriginCell": {
			MouseAt: func() pixel.Vec { return pixel.V(-grid.CellHeight/2, 0) },
		},
		"AtRightEdgeOfOriginCell": {
			MouseAt: func() pixel.Vec { return pixel.V(grid.CellHeight/2, 0) },
		},
		"AtTopEdgeOfOriginCell": {
			MouseAt: func() pixel.Vec { return pixel.V(0, grid.CellHeight/2) },
		},
		"AtBottomEdgeOfOriginCell": {
			MouseAt: func() pixel.Vec { return pixel.V(0, -grid.CellHeight/2) },
		},

		"AtMiddleOfCellRightOfOrigin": {
			MouseAt: func() pixel.Vec { return pixel.V(grid.CellWidth, 0) },
			Cell:    game.P(0, 1),
		},
		"AtMiddleOfCellLeftOfOrigin": {
			MouseAt: func() pixel.Vec { return pixel.V(-grid.CellWidth, 0) },
			Cell:    game.P(0, -1),
		},
		"AtMiddleOfCellAboveOrigin": {
			MouseAt: func() pixel.Vec { return pixel.V(0, grid.CellHeight) },
			Cell:    game.P(1, 0),
		},
		"AtMiddleOfCellBellowOrigin": {
			MouseAt: func() pixel.Vec { return pixel.V(0, -grid.CellHeight) },
			Cell:    game.P(-1, 0),
		},

		"LookingAtLeftEdgeOfOriginCell": {
			LookingAt: pixel.V(-grid.CellHeight/2, 0),
		},
		"LookingAtRightEdgeOfOriginCell": {
			LookingAt: pixel.V(grid.CellHeight/2, 0),
		},
		"LookingAtTopEdgeOfOriginCell": {
			LookingAt: pixel.V(0, grid.CellHeight/2),
		},
		"LookingAtBottomEdgeOfOriginCell": {
			LookingAt: pixel.V(0, -grid.CellHeight/2),
		},

		"LookingAtMiddleOfCellRightOfOrigin": {
			LookingAt: pixel.V(grid.CellWidth, 0),
			Cell:      game.P(0, 1),
		},
		"LookingAtMiddleOfCellLeftOfOrigin": {
			LookingAt: pixel.V(-grid.CellWidth, 0),
			Cell:      game.P(0, -1),
		},
		"LookingAtMiddleOfCellAboveOrigin": {
			LookingAt: pixel.V(0, grid.CellHeight),
			Cell:      game.P(1, 0),
		},
		"LookingAtMiddleOfCellBellowOrigin": {
			LookingAt: pixel.V(0, -grid.CellHeight),
			Cell:      game.P(-1, 0),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			input := TestInput{}
			input.Mock.MousePosition = tt.MouseAt

			cam := ui.NewCamera(tt.LookingAt)

			// when
			cell := grid.UnderCursor(input, cam)

			// then
			assert.That(cell == tt.Cell, t.Errorf, "got cell %#v, want %#v", cell, tt.Cell)
		})
	}
}
