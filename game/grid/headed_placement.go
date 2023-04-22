// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package grid

import (
	"time"

	"github.com/szabba/tob-cob/game/actions"
)

// A HeadedPlacement is either:
//
//   - nowhere
//   - at a point
//   - moving from one point to another
type HeadedPlacement struct {
	pos, heading OnePosTaker
	countdown    actions.Countdown
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
func (hp *HeadedPlacement) MoveTo(dst Position, dt time.Duration) actions.Action {
	if !hp.Placed() {
		return actions.NoAction()
	}
	return actions.Sequence(
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
	return hp.countdown.Progress()
}

func (hp *HeadedPlacement) arriveAction(dst Position) actions.Action {
	return &_PlacemnetArriveAction{hp, dst}
}

type _PlacemnetArriveAction struct {
	placement *HeadedPlacement
	dst       Position
}

func (action *_PlacemnetArriveAction) Run(atMost time.Duration) actions.Status {
	action.placement.heading.Leave()
	// We just made dst available - therefore Take will succeed.
	action.dst.Take(&action.placement.pos)
	return actions.Done(atMost)
}

// FollowPath creates an action that will move the heading along the path.
// Each step will take stepDt.
//
// When first run, the action fails immediately if the heading is not at the initial position of the path.
func (hp *HeadedPlacement) FollowPath(path Path, stepDt time.Duration) actions.Action {
	steps := make([]actions.Action, 0, len(path))
	if len(path) > 0 {
		steps = append(steps, hp.checkAt(path[0]))
		for _, dst := range path[1:] {
			steps = append(steps, hp.MoveTo(dst, stepDt))
		}
	}
	return actions.Sequence(steps...)
}

func (hp *HeadedPlacement) checkAt(pos Position) actions.Action {
	return &_PlacementCheckAtAction{hp, pos}
}

type _PlacementCheckAtAction struct {
	placement *HeadedPlacement
	pos       Position
}

func (action *_PlacementCheckAtAction) Run(atMost time.Duration) actions.Status {
	if action.placement.pos.pos != action.pos {
		return actions.Interrupted(atMost)
	}
	return actions.Done(atMost)
}
