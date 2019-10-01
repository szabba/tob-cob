// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

import "time"

// An Action is a process stretched out in time.
type Action interface {
	// Run runs the action, possibly to completion.
	// The parameter says for how long the action is allowed to run.
	//
	// The returned status says if the action has completed.
	// It also says how much time is left after it has run.
	Run(timeLeft time.Duration) ActionStatus
}

// An ActionStatus says whether an action completed.
// It also says how much time is left after running it.
type ActionStatus struct {
	timeLeft time.Duration
	done     bool
}

// HasTimeLeft says whether there is time left after running an action.
func (status ActionStatus) HasTimeLeft() bool { return status.TimeLeft() > 0 }

// TimeLeft returns the time still left after running the action.
func (status ActionStatus) TimeLeft() time.Duration { return status.timeLeft }

// Done says whether the action has completed.
func (status ActionStatus) Done() bool { return status.done }

// Done creates an action status indicating that an action has completed.
//
// The timeLeft should be the time the action did not use up.
// As a special case, there will be no time left when the argument is negative.
func Done(timeLeft time.Duration) ActionStatus {
	if timeLeft < 0 {
		timeLeft = 0
	}
	return ActionStatus{timeLeft: timeLeft, done: true}
}

// Paused creates an action status indicating that an action ran out of time.
//
// This means that the action has not completed.
// Also, there is no time left to run other actions.
func Paused() ActionStatus {
	return ActionStatus{timeLeft: 0, done: false}
}

// NoAction returns an action that completes instantly.
//
// Running the returned action has no side-effects.
func NoAction() Action {
	return _noAction
}

type _NoAction struct{}

var _noAction = &_NoAction{}

func (*_NoAction) Run(timeLeft time.Duration) ActionStatus {
	return Done(timeLeft)
}
