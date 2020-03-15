// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package actions

import (
	"time"
)

// An Action is a process stretched out in time.
type Action interface {
	// Run runs the action, possibly to completion.
	// The parameter says for how long the action is allowed to run.
	//
	// The returned status says if the action has completed.
	// It also says how much time is left after it has run.
	Run(atMost time.Duration) Status
}

// An Status says whether an action completed.
// It also says how much time is left after running it.
type Status struct {
	timeLeft    time.Duration
	done        bool
	interrupted bool
}

// Interrupted says whether the action was interrupted.
func (status Status) Interrupted() bool { return status.interrupted }

// HasTimeLeft says whether there is time left after running an action.
func (status Status) HasTimeLeft() bool { return status.TimeLeft() > 0 }

// TimeLeft returns the time still left after running the action.
func (status Status) TimeLeft() time.Duration { return status.timeLeft }

// Done says whether the action has completed.
func (status Status) Done() bool { return status.done }

// Interrupted creates an action status saying that the action was interrupted.
//
// The timeLeft should be the time the action did not use up.
// As a special case, there will be no time left when the argument is negative.
func Interrupted(timeLeft time.Duration) Status {
	if timeLeft < 0 {
		timeLeft = 0
	}
	return Status{timeLeft: timeLeft, interrupted: true}
}

// Done creates an action status indicating that an action has completed.
//
// The timeLeft should be the time the action did not use up.
// As a special case, there will be no time left when the argument is negative.
func Done(timeLeft time.Duration) Status {
	if timeLeft < 0 {
		timeLeft = 0
	}
	return Status{timeLeft: timeLeft, done: true}
}

// Paused creates an action status indicating that an action needs more time to complete.
func Paused() Status {
	return Status{}
}

// NoAction returns an action that completes instantly.
//
// Running the returned action has no side-effects.
func NoAction() Action {
	return _noAction
}

type _NoAction struct{}

var _noAction = &_NoAction{}

func (*_NoAction) Run(atMost time.Duration) Status {
	return Done(atMost)
}

// Interrupt returns an action that immediately interrupts.
func Interrupt() Action { return _interrupt }

type _Interrupt struct{}

var _interrupt = &_Interrupt{}

func (*_Interrupt) Run(atMost time.Duration) Status {
	return Interrupted(atMost)
}
