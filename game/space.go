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
	_, exists := pos.space.poses[pos.at]
	return exists
}

// Create ensures that a position within a space exists.
// It fails when the position exists.
func (pos Position) Create() bool {
	if pos.Exists() {
		return false
	}
	pos.space.poses[pos.at] = pos.Taken()
	return true
}

// Destroy ensures that a position within a space does not exist.
// It fails when the position does not exist or is taken.
func (pos Position) Destroy() bool {
	if pos.Taken() {
		return false
	}
	ok := pos.Exists()
	delete(pos.space.poses, pos.at)
	return ok
}

// Taken says whether the position is currently taken.
func (pos Position) Taken() bool {
	return pos.space.poses[pos.at]
}

// Take tries to mark the position as taken.
// It fails if the position does not exist or is free.
func (pos Position) Take() bool {
	if !pos.Exists() || pos.Taken() {
		return false
	}
	pos.space.poses[pos.at] = true
	return true
}

// Free tries to mark the position as no longer taken.
// It fails if the position is not taken.
func (pos Position) Free() bool {
	if !pos.Taken() {
		return false
	}
	pos.space.poses[pos.at] = false
	return true
}
