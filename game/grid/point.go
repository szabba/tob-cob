// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package grid

// A Point is a point on a 2D grid.
type Point struct {
	Row, Column int
}

// P creates a position at the said row and column.
func P(row, column int) Point {
	return Point{Row: row, Column: column}
}
