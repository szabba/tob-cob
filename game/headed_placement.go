// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

import (
	"time"
)

// A HeadedPlacement is either:
//
//     - nowhere
//     - at a point
//     - moving from one point to another
type HeadedPlacement struct {
	pos, heading OnePosTaker
	countdown    Countdown
}

// Placed says whether the placement takes up some position.
// Nothing else is be able to occupy the same position at the same time.
func (hp *HeadedPlacement) Placed() bool { return hp.pos.Placed() }

// AtPoint the current point of current the placement.
// As a special case it returns the zero value when the placement is not placed anywhere.
// Call Placed to distinguish that from a placement at that position.
func (hp *HeadedPlacement) AtPoint() Point {
	return hp.pos.AtPoint()
}

// Headed says whether the placement is headed somewhere.
// Nothing else is able to occupy the position where the placement is headed at the same time.
func (hp *HeadedPlacement) Headed() bool {
	return hp.heading.Placed()
}

// Heading is the point where the placement is headed.
// As a special case it returns the zero value when the placement is headed nowhere.
// Call Headed to distinguish that from a placement headed there.
func (hp *HeadedPlacement) Heading() Point {
	return hp.heading.AtPoint()
}

// Place tries to put the placement at the given position.
// The return value says whether this has succeded.
func (hp *HeadedPlacement) Place(pos Position) bool {
	pos.Take(&hp.pos)
	hp.heading.Leave()
	return hp.Placed()
}

// MoveTo creates an action that will try to move the placement to dst over a duration of dt.
// The action gets interrupted if the dst position is taken as it starts.
func (hp *HeadedPlacement) MoveTo(dst Position, dt time.Duration) Action {
	if !hp.Placed() {
		return NoAction()
	}
	return Sequence(
		TakePosition(dst, &hp.heading),
		hp.countdown.Action(dt),
		hp.arriveAction(dst))
}

// Progress says how far along the placement is in the move from the start position to the heading.
// This is 1 when the placement is not headed anywhere.
func (hp *HeadedPlacement) Progress() float64 {
	if !hp.Placed() {
		return 0
	}
	if hp.Placed() && !hp.Headed() {
		return 1
	}
	return hp.countdown.Progress()
}

func (hp *HeadedPlacement) arriveAction(dst Position) Action {
	return &_PlacemnetArriveAction{hp, dst}
}

type _PlacemnetArriveAction struct {
	placement *HeadedPlacement
	dst       Position
}

func (action *_PlacemnetArriveAction) Run(atMost time.Duration) ActionStatus {
	action.placement.heading.Leave()
	action.dst.Take(&action.placement.pos)
	return Done(atMost)
}
