// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package actions

import (
	"time"
)

// Sequence creates an action that runs several steps one after another.
func Sequence(steps ...Action) Action {
	return &_Sequence{steps}
}

type _Sequence struct {
	steps []Action
}

func (seq *_Sequence) Run(atMost time.Duration) Status {
	status := Done(atMost)
	for status.Done() && seq.hasStepsLeft() {
		status = seq.runStep(status.TimeLeft())
	}
	if !seq.hasStepsLeft() {
		return status
	}
	return Paused()
}

func (seq *_Sequence) hasStepsLeft() bool {
	return len(seq.steps) > 0
}

func (seq *_Sequence) runStep(atMost time.Duration) Status {
	status := seq.steps[0].Run(atMost)
	if status.Interrupted() {
		seq.steps = nil
	} else if status.Done() {
		seq.steps = seq.steps[1:]
	}
	return status
}
