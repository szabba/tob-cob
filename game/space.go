// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

import "time"

// A Space where things can exist and interact.
//
// It is a subspace of a 2D grid.
// Which positions on the grid exist can change dynamically.
type Space struct {
	poses map[Point]SpaceTaker
}

// NewSpace creates a new, empty space.
func NewSpace() *Space {
	return &Space{
		poses: map[Point]SpaceTaker{},
	}
}

// At returns the position at the given point in the space.
//
// Two position values returned for the same point in the space will be equal.
func (space *Space) At(at Point) Position {
	return Position{space, at}
}

// A Position within some space.
type Position struct {
	space *Space
	at    Point
}

// AtPoint is the point at which the position is in the space that contains it.
func (pos Position) AtPoint() Point { return pos.at }

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
	pos.space.poses[pos.at] = nil
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
	return pos.space.poses[pos.at] != nil
}

// Take tries to mark the position as taken.
// It fails if the position does not exist or is free.
func (pos Position) Take(taker SpaceTaker) bool {
	if !pos.Exists() || pos.Taken() {
		return false
	}
	taker.LetOnto(pos)
	pos.space.poses[pos.at] = taker
	return true
}

// Free tries to mark the position as no longer taken.
// It fails if the position is not taken.
func (pos Position) Free() bool {
	if !pos.Taken() {
		return false
	}
	pos.space.poses[pos.at].ForceOff(pos)
	pos.space.poses[pos.at] = nil
	return true
}

// A SpaceTaker is the thing that takes up a taken position.
//
// A space taker can take up multiple positions at once.
type SpaceTaker interface {
	// LetOnto tells the space taker that it is now taking up pos.
	LetOnto(pos Position)
	// ForceOff tells the space taker that it is no longer taking up pos.
	ForceOff(pos Position)
}

// DummyTaker returns a space taker that only takes up the position it takes.
// It has no additional behaviour.
func DummyTaker() SpaceTaker {
	return _dummyTaker
}

type _DummyTaker struct{}

var _dummyTaker = &_DummyTaker{}

func (*_DummyTaker) LetOnto(_ Position)  {}
func (*_DummyTaker) ForceOff(_ Position) {}

// TakePositionAction is an action that tries to immediately take a position.
type TakePositionAction struct {
	pos   Position
	taker SpaceTaker
}

// TakePosition builds a TakePositionAcion.
func TakePosition(pos Position, taker SpaceTaker) *TakePositionAction {
	return &TakePositionAction{
		pos:   pos,
		taker: taker,
	}
}

var _ Action = &TakePositionAction{}

// Run immediately takes the position or fails interrupting the action.
func (action *TakePositionAction) Run(atMost time.Duration) ActionStatus {
	if action.pos.Take(action.taker) {
		return Done(atMost)
	}
	return Interrupted(atMost)
}
