// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui_test

import (
	"testing"

	"github.com/szabba/assert"
	"github.com/szabba/tob-cob/game/grid"
	"github.com/szabba/tob-cob/ui"
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input/testinput"
)

func TestGridDimmensionsMatrix(t *testing.T) {
	dims := ui.GridDimensions{CellWidth: 20, CellHeight: 10}
	tests := map[string]struct {
		Column, Row int

		RelativeToCell  geometry.Vec
		RelativeToWorld geometry.Vec
	}{
		"MiddleOfOriginCell": {},

		"TopEdgeOfOriginCell": {
			RelativeToCell:  geometry.V(0, dims.CellHeight/2),
			RelativeToWorld: geometry.V(0, dims.CellHeight/2),
		},
		"BottomEdgeOfOriginCell": {
			RelativeToCell:  geometry.V(0, -dims.CellHeight/2),
			RelativeToWorld: geometry.V(0, -dims.CellHeight/2),
		},
		"LeftEdgeOfOriginCell": {
			RelativeToCell:  geometry.V(-dims.CellWidth/2, 0),
			RelativeToWorld: geometry.V(-dims.CellWidth/2, 0),
		},
		"RightEdgeOfOriginCell": {
			RelativeToCell:  geometry.V(dims.CellWidth/2, 0),
			RelativeToWorld: geometry.V(dims.CellWidth/2, 0),
		},

		"MiddleOfCellAboveOrigin": {
			Row:             1,
			RelativeToWorld: geometry.V(0, dims.CellHeight),
		},
		"MiddleOfCellBelowOrigin": {
			Row:             -1,
			RelativeToWorld: geometry.V(0, -dims.CellHeight),
		},
		"MiddleOfCellRightOfOrigin": {
			Column:          1,
			RelativeToWorld: geometry.V(dims.CellWidth, 0),
		},
		"MiddleOfCellLeftOfOrigin": {
			Column:          -1,
			RelativeToWorld: geometry.V(-dims.CellWidth, 0),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			matrix := dims.Matrix(tt.Column, tt.Row)

			// then
			relativeToWorld := matrix.Apply(tt.RelativeToCell)
			assert.That(
				relativeToWorld == tt.RelativeToWorld,
				t.Errorf, "got %#v world-coordinate, want %#v", relativeToWorld, tt.RelativeToWorld)
		})
	}
}

func TestGridDimmensionsUnderCursor(t *testing.T) {
	dims := ui.GridDimensions{CellWidth: 20, CellHeight: 10}
	tests := map[string]struct {
		LookingAt geometry.Vec
		MouseAt   func() geometry.Vec

		Cell grid.Point
	}{
		"AtOrigin": {},

		"AtLeftEdgeOfOriginCell": {
			MouseAt: func() geometry.Vec { return geometry.V(-dims.CellHeight/2, 0) },
		},
		"AtRightEdgeOfOriginCell": {
			MouseAt: func() geometry.Vec { return geometry.V(dims.CellHeight/2, 0) },
		},
		"AtTopEdgeOfOriginCell": {
			MouseAt: func() geometry.Vec { return geometry.V(0, dims.CellHeight/2) },
		},
		"AtBottomEdgeOfOriginCell": {
			MouseAt: func() geometry.Vec { return geometry.V(0, -dims.CellHeight/2) },
		},

		"AtMiddleOfCellRightOfOrigin": {
			MouseAt: func() geometry.Vec { return geometry.V(dims.CellWidth, 0) },
			Cell:    grid.P(0, 1),
		},
		"AtMiddleOfCellLeftOfOrigin": {
			MouseAt: func() geometry.Vec { return geometry.V(-dims.CellWidth, 0) },
			Cell:    grid.P(0, -1),
		},
		"AtMiddleOfCellAboveOrigin": {
			MouseAt: func() geometry.Vec { return geometry.V(0, dims.CellHeight) },
			Cell:    grid.P(1, 0),
		},
		"AtMiddleOfCellBellowOrigin": {
			MouseAt: func() geometry.Vec { return geometry.V(0, -dims.CellHeight) },
			Cell:    grid.P(-1, 0),
		},

		"LookingAtLeftEdgeOfOriginCell": {
			LookingAt: geometry.V(-dims.CellHeight/2, 0),
		},
		"LookingAtRightEdgeOfOriginCell": {
			LookingAt: geometry.V(dims.CellHeight/2, 0),
		},
		"LookingAtTopEdgeOfOriginCell": {
			LookingAt: geometry.V(0, dims.CellHeight/2),
		},
		"LookingAtBottomEdgeOfOriginCell": {
			LookingAt: geometry.V(0, -dims.CellHeight/2),
		},

		"LookingAtMiddleOfCellRightOfOrigin": {
			LookingAt: geometry.V(dims.CellWidth, 0),
			Cell:      grid.P(0, 1),
		},
		"LookingAtMiddleOfCellLeftOfOrigin": {
			LookingAt: geometry.V(-dims.CellWidth, 0),
			Cell:      grid.P(0, -1),
		},
		"LookingAtMiddleOfCellAboveOrigin": {
			LookingAt: geometry.V(0, dims.CellHeight),
			Cell:      grid.P(1, 0),
		},
		"LookingAtMiddleOfCellBellowOrigin": {
			LookingAt: geometry.V(0, -dims.CellHeight),
			Cell:      grid.P(-1, 0),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			input := testinput.Source{}
			input.Mock.MousePosition = tt.MouseAt

			cam := ui.NewCamera(tt.LookingAt)

			// when
			cell := dims.UnderCursor(input, cam)

			// then
			assert.That(cell == tt.Cell, t.Errorf, "got cell %#v, want %#v", cell, tt.Cell)
		})
	}
}
