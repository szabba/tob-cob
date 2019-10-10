// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

// A Space where things can exist and interact.
//
// It is a subspace of a 2D grid.
// Which positions on the grid exist can change dynamically.
type Space struct {
	poses map[Point]bool
}

// NewSpace creates a new, empty space.
func NewSpace() *Space {
	return &Space{
		poses: map[Point]bool{},
	}
}

// At returns the position at the given point in the space.
func (space *Space) At(at Point) Position {
	return Position{space, at}
}

// A Position within some space.
type Position struct {
	space *Space
	at    Point
}

// Exists says whether the position within the space exists.
func (pos Position) Exists() bool {
	return pos.space.poses[pos.at]
}

// Create ensures that a position within a space exists.
func (pos Position) Create() {
	pos.space.poses[pos.at] = true
}

// Destroy ensures that a position within a space does not exist.
func (pos Position) Destroy() {
	delete(pos.space.poses, pos.at)
}
