// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

import (
	"sync"
	"time"
)

// A Countdown tracks progress towards reaching a numerical goal.
// The zero value is a completed countdown with a goal of 0.
type Countdown struct {
	elapsed, needed int
}

// ResetTarget sets a new target value and resets all progress made.
func (c *Countdown) ResetTarget(needed int) {
	c.elapsed = 0
	c.needed = needed
}

// Progress says how far along the countdown is.
// The value ranges from 0 (not started) to 1 (completed).
func (c Countdown) Progress() float64 {
	if c.needed == 0 {
		return 1
	}
	return float64(c.elapsed) / float64(c.needed)
}

// CountDown tracks progress being made towards a countdown.
// The countdown gets closer to completion with bigger arguments.
//
func (c *Countdown) CountDown(atMost int) (leftOver int) {
	neededLeft := c.needed - c.elapsed
	if atMost >= neededLeft {
		c.elapsed = c.needed
		return atMost - neededLeft
	}
	c.elapsed += atMost
	return 0
}

// Action creates an action that makes the countdown progress run.
// The countdown completes when the action does.
func (c *Countdown) Action(lasting time.Duration) Action {
	return &_CountdownAction{lasting: lasting, countdown: c}
}

type _CountdownAction struct {
	once      sync.Once
	lasting   time.Duration
	countdown *Countdown
}

func (action *_CountdownAction) Run(atMost time.Duration) ActionStatus {
	action.once.Do(action.init)
	leftOver := action.countdown.CountDown(int(atMost))
	if action.countdown.Progress() < 1 {
		return Paused()
	}
	return Done(time.Duration(leftOver))
}

func (action *_CountdownAction) init() {
	needed := int(action.lasting)
	action.countdown.ResetTarget(needed)
}
