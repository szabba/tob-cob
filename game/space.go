// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

import (
	"time"

	"github.com/szabba/tob-cob/game/actions"
)

// A Space where things can exist and interact.
//
// It is a subspace of a 2D grid.
// Which positions on the grid exist can change dynamically.
type Space struct {
	poses    map[Point]SpaceTaker
	min, max Point
	empty    bool
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

// Min is the point with the column of the leftmost position and row of the bottom one.
func (space *Space) Min() Point { return space.min }

// Max is the point with the column of the rightmost position and row of the top one.
func (space *Space) Max() Point { return space.max }

func (space *Space) fixMinMax() {
	space.empty = true
	space.min, space.max = Point{}, Point{}
	for at := range space.poses {
		space.pickMin(at.Row, &space.min.Row)
		space.pickMin(at.Column, &space.min.Column)
		space.pickMax(at.Row, &space.max.Row)
		space.pickMax(at.Column, &space.max.Column)
		space.empty = false
	}
}

func (space *Space) pickMin(candidate int, dst *int) {
	if space.empty || candidate < *dst {
		*dst = candidate
	}
}

func (space *Space) pickMax(candidate int, dst *int) {
	if space.empty || *dst < candidate {
		*dst = candidate
	}
}

// A Position within some space.
type Position struct {
	space *Space
	at    Point
}

// AtPoint is the point at which the position is in the space that contains it.7
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
	pos.space.fixMinMax()
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
	pos.space.fixMinMax()
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

// TakePosition builds an action that tries to immediately take the given position with the given taker.
// If the position is taken, the action fails, interrupted.
func TakePosition(pos Position, taker SpaceTaker) actions.Action {
	return _TakePositionAction{
		pos:   pos,
		taker: taker,
	}
}

type _TakePositionAction struct {
	pos   Position
	taker SpaceTaker
}

var _ actions.Action = _TakePositionAction{}

func (action _TakePositionAction) Run(atMost time.Duration) actions.Status {
	if action.pos.Take(action.taker) {
		return actions.Done(atMost)
	}
	return actions.Interrupted(atMost)
}
