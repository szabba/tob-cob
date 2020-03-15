// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package actions_test

import (
	"testing"
	"time"

	"github.com/szabba/assert"
	"github.com/szabba/tob-cob/game/actions"
)

func TestSequence(t *testing.T) {
	tests := map[string]struct {
		Steps  []actions.Action
		Times  []time.Duration
		Status actions.Status
	}{
		"Empty": {
			Times:  []time.Duration{time.Second},
			Status: actions.Done(time.Second),
		},

		"SingleStep/InsufficientTime": {
			Steps:  []actions.Action{actions.Wait(3 * time.Second)},
			Times:  []time.Duration{2 * time.Second},
			Status: actions.Paused(),
		},
		"SingleStep/ExactTime": {
			Steps:  []actions.Action{actions.Wait(time.Second)},
			Times:  []time.Duration{time.Second},
			Status: actions.Done(0),
		},
		"SingleStep/TimeToSpare": {
			Steps:  []actions.Action{actions.Wait(time.Second)},
			Times:  []time.Duration{3 * time.Second},
			Status: actions.Done(2 * time.Second),
		},

		"TwoSteps/First/InsufficientTime": {
			Steps: []actions.Action{
				actions.Wait(2 * time.Second),
				actions.Wait(3 * time.Second),
			},
			Times:  []time.Duration{time.Second},
			Status: actions.Paused(),
		},
		"TwoSteps/First/ExactTime": {
			Steps: []actions.Action{
				actions.Wait(2 * time.Second),
				actions.Wait(3 * time.Second),
			},
			Times:  []time.Duration{2 * time.Second},
			Status: actions.Paused(),
		},
		"TwoSteps/Second/InsufficientTime": {
			Steps: []actions.Action{
				actions.Wait(2 * time.Second),
				actions.Wait(3 * time.Second),
			},
			Times:  []time.Duration{4 * time.Second},
			Status: actions.Paused(),
		},
		"TwoSteps/Second/ExactTime": {
			Steps: []actions.Action{
				actions.Wait(2 * time.Second),
				actions.Wait(3 * time.Second),
			},
			Times:  []time.Duration{5 * time.Second},
			Status: actions.Done(0),
		},
		"TwoSteps/Second/TimeToSpare": {
			Steps: []actions.Action{
				actions.Wait(2 * time.Second),
				actions.Wait(3 * time.Second),
			},
			Times:  []time.Duration{6 * time.Second},
			Status: actions.Done(time.Second),
		},

		"Interrupted": {
			Steps: []actions.Action{
				actions.Interrupt(),
				actions.Wait(time.Second),
			},
			Times:  []time.Duration{2 * time.Second},
			Status: actions.Interrupted(2 * time.Second),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			assert.That(len(tt.Times) > 0, t.Errorf, "case %q has no times specified", name)
			action := actions.Sequence(tt.Steps...)

			// when
			var status actions.Status
			for _, t := range tt.Times {
				status = action.Run(t)
			}

			// then
			assert.That(
				status == tt.Status,
				t.Errorf, "final action status is %#v, want %#v", status, tt.Status)
		})
	}
}

func TestSequenceRunsFirstActionForZeroTimeWhenNoneIsGiven(t *testing.T) {
	// given
	action := &StepRecordingAction{}
	sequence := actions.Sequence(action)

	// when
	sequence.Run(0)

	// then
	assert.That(len(action.Steps) == 1, t.Fatalf, "got %d steps run - want %d", len(action.Steps), 1)
	assert.That(action.Steps[0] == 0, t.Errorf, "got step %d of length %s - want %s", 0, action.Steps[0], 0)
}

func TestSequenceRunsRemainingStepForZeroTimeIfNoneIsLeft(t *testing.T) {
	// given
	action := &StepRecordingAction{}
	sequence := actions.Sequence(actions.Wait(time.Second), action)

	// when
	sequence.Run(time.Second)

	// then
	assert.That(len(action.Steps) == 1, t.Fatalf, "got %d steps run - want %d", len(action.Steps), 1)
	assert.That(action.Steps[0] == 0, t.Errorf, "got step %d of length %s - want %s", 0, action.Steps[0], 0)
}

type StepRecordingAction struct {
	Steps []time.Duration
}

func (action *StepRecordingAction) Run(atMost time.Duration) actions.Status {
	action.Steps = append(action.Steps, atMost)
	return actions.Paused()
}
