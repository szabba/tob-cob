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
	timeLeft    time.Duration
	done        bool
	interrupted bool
}

// Interrupted says whether the action was interrupted.
func (status ActionStatus) Interrupted() bool { return status.interrupted }

// HasTimeLeft says whether there is time left after running an action.
func (status ActionStatus) HasTimeLeft() bool { return status.TimeLeft() > 0 }

// TimeLeft returns the time still left after running the action.
func (status ActionStatus) TimeLeft() time.Duration { return status.timeLeft }

// Done says whether the action has completed.
func (status ActionStatus) Done() bool { return status.done }

// Interrupted creates an action status saying that the action was interrupted.
//
// The timeLeft should be the time the action did not use up.
// As a special case, there will be no time left when the argument is negative.
func Interrupted(timeLeft time.Duration) ActionStatus {
	if timeLeft < 0 {
		timeLeft = 0
	}
	return ActionStatus{timeLeft: timeLeft, interrupted: true}
}

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

// Paused creates an action status indicating that an action needs more time to complete.
func Paused() ActionStatus {
	return ActionStatus{}
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

// Interrupt returns an action that immediately interrupts.
func Interrupt() Action { return _interrupt }

type _Interrupt struct{}

var _interrupt = &_Interrupt{}

func (*_Interrupt) Run(timeLeft time.Duration) ActionStatus {
	return Interrupted(timeLeft)
}

// Wait returns an action that lasts waitTime but does nothing.
func Wait(waitTime time.Duration) Action {
	return &_Wait{waitTime}
}

type _Wait struct {
	toEnd time.Duration
}

func (w *_Wait) Run(timeLeft time.Duration) ActionStatus {
	if timeLeft < w.toEnd {
		w.toEnd -= timeLeft
		return Paused()
	}
	return Done(timeLeft - w.toEnd)
}

// Sequence creates an action that runs several steps one after another.
func Sequence(steps ...Action) Action {
	return &_Sequence{steps}
}

type _Sequence struct {
	steps []Action
}

func (seq *_Sequence) Run(timeLeft time.Duration) ActionStatus {
	lastStatus := Done(timeLeft)
	for lastStatus.HasTimeLeft() && seq.hasStepsLeft() {
		lastStatus = seq.runStep(lastStatus.TimeLeft())
	}
	if !seq.hasStepsLeft() {
		return lastStatus
	}
	return Paused()
}

func (seq *_Sequence) hasStepsLeft() bool {
	return len(seq.steps) > 0
}

func (seq *_Sequence) runStep(timeLeft time.Duration) ActionStatus {
	status := seq.steps[0].Run(timeLeft)
	if status.Interrupted() {
		seq.steps = nil
	} else if status.Done() {
		seq.steps = seq.steps[1:]
	}
	return status
}
