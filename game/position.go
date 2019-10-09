// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

// A Position is a point on a 2D grid.
type Position struct {
	Row, Column int
}

// P creates a position at the said row and column.
func P(row, column int) Position {
	return Position{Row: row, Column: column}
}
