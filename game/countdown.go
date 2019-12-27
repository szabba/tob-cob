// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

import "time"

// A Countdown tracks progress towards reaching a numerical goal.
type Countdown struct {
	elapsed, needed int
}

// CountdownTo creates a countdown with the specified targer value.
func CountdownTo(needed int) Countdown {
	return Countdown{needed: needed}
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

// CountdownOver creates a countdown with an associated action.
// The countdown completes when the action does.
func CountdownOver(needed time.Duration) (*Countdown, Action) {
	progress := CountdownTo(int(needed))
	return &progress, &_ProgressAction{&progress}
}

type _ProgressAction struct {
	progress *Countdown
}

func (action *_ProgressAction) Run(atMost time.Duration) ActionStatus {
	leftOver := action.progress.CountDown(int(atMost))
	if action.progress.Progress() < 1 {
		return Paused()
	}
	return Done(time.Duration(leftOver))
}
