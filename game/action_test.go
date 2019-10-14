// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game_test

import (
	"testing"
	"time"

	"github.com/szabba/assert"
	"github.com/szabba/tob-cob/game"
)

func TestInterrupted(t *testing.T) {
	// given
	timeLeft := time.Second

	// when
	status := game.Interrupted(timeLeft)

	// then
	assert.That(!status.Done(), t.Errorf, "the status should say the action has not completed")
	assert.That(
		status.Interrupted(),
		t.Errorf, "the status should say the action was interrupted")
	assert.That(
		status.HasTimeLeft(),
		t.Errorf, "the status should say there is time left after running the action")
	assert.That(
		status.TimeLeft() == timeLeft,
		t.Errorf, "time left is %s, want %s", status.TimeLeft(), timeLeft)
}

func TestInterruptedWithNoTimeLeft(t *testing.T) {
	// given
	timeLeft := time.Duration(0)

	// when
	status := game.Interrupted(timeLeft)

	// then
	assert.That(!status.Done(), t.Errorf, "the status should say the action has not completed")
	assert.That(
		status.Interrupted(),
		t.Errorf, "the status should say the action was interrupted")
	assert.That(
		!status.HasTimeLeft(),
		t.Errorf, "the status should say there is no time left after running the action")
	assert.That(
		status.TimeLeft() == timeLeft,
		t.Errorf, "time left is %s, want %s", status.TimeLeft(), timeLeft)
}

func TestInterruptedWithNegativeTimeLeft(t *testing.T) {
	// given
	timeLeft := -time.Second
	sanitizedTime := time.Duration(0)

	// when
	status := game.Interrupted(timeLeft)

	// then
	assert.That(!status.Done(), t.Errorf, "the status should say the action has not completed")
	assert.That(
		status.Interrupted(),
		t.Errorf, "the status should say the action was interrupted")
	assert.That(
		!status.HasTimeLeft(),
		t.Errorf, "the status should say there is no time left after running the action")
	assert.That(
		status.TimeLeft() == sanitizedTime,
		t.Errorf, "time left is %s, want %s", status.TimeLeft(), sanitizedTime)
}
func TestDone(t *testing.T) {
	// given
	timeLeft := time.Second

	// when
	status := game.Done(timeLeft)

	// then
	assert.That(status.Done(), t.Errorf, "the status should say the action has completed")
	assert.That(
		!status.Interrupted(),
		t.Errorf, "the status should say the action was not interrupted")
	assert.That(
		status.HasTimeLeft(),
		t.Errorf, "the status should say there is time left after running the action")
	assert.That(
		status.TimeLeft() == timeLeft,
		t.Errorf, "time left is %s, want %s", status.TimeLeft(), timeLeft)
}

func TestDoneWithNoTimeLeft(t *testing.T) {
	// given
	timeLeft := time.Duration(0)

	// when
	status := game.Done(timeLeft)

	// then
	assert.That(status.Done(), t.Errorf, "the status should say the action has completed")
	assert.That(
		!status.Interrupted(),
		t.Errorf, "the status should say the action was not interrupted")
	assert.That(
		!status.HasTimeLeft(),
		t.Errorf, "the status should say there is no time left after running the action")
	assert.That(
		status.TimeLeft() == timeLeft,
		t.Errorf, "time left is %s, want %s", status.TimeLeft(), timeLeft)
}

func TestDoneWithNegativeTimeLeft(t *testing.T) {
	// given
	timeLeft := -time.Second
	sanitizedTime := time.Duration(0)

	// when
	status := game.Done(timeLeft)

	// then
	assert.That(status.Done(), t.Errorf, "the status should say the action has completed")
	assert.That(
		!status.Interrupted(),
		t.Errorf, "the status should say the action was not interrupted")
	assert.That(
		!status.HasTimeLeft(),
		t.Errorf, "the status should say there is no time left after running the action")
	assert.That(
		status.TimeLeft() == sanitizedTime,
		t.Errorf, "time left is %s, want %s", status.TimeLeft(), sanitizedTime)
}

func TestPause(t *testing.T) {
	// given
	expectedTime := time.Duration(0)

	// when
	status := game.Paused()

	// then
	// then
	assert.That(!status.Done(), t.Errorf, "the status should say the action has not completed")
	assert.That(
		!status.Interrupted(),
		t.Errorf, "the status should say the action was not interrupted")
	assert.That(
		!status.HasTimeLeft(),
		t.Errorf, "the status should say there is no time left after running the action")
	assert.That(
		status.TimeLeft() == expectedTime,
		t.Errorf, "time left is %s, want %s", status.TimeLeft(), expectedTime)
}

func TestNoActionCompletesInNoTime(t *testing.T) {
	// given
	timeLeft := time.Second
	action := game.NoAction()
	statusWanted := game.Done(timeLeft)

	// when
	status := action.Run(timeLeft)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}

func TestInterruptsFailsInNoTime(t *testing.T) {
	// given
	timeLeft := time.Second
	action := game.Interrupt()
	statusWanted := game.Interrupted(timeLeft)

	// when
	status := action.Run(timeLeft)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}

func TestWaitWithExactTime(t *testing.T) {
	// given
	action := game.Wait(time.Second)
	statusWanted := game.Done(0)

	// when
	status := action.Run(time.Second)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}

func TestWaitWithoutEnoughTime(t *testing.T) {
	// given
	action := game.Wait(3 * time.Second)
	statusWanted := game.Paused()

	// when
	status := action.Run(time.Second)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}
func TestWaitWithTimeToSpare(t *testing.T) {
	// given
	action := game.Wait(time.Second)
	statusWanted := game.Done(2 * time.Second)

	// when
	status := action.Run(3 * time.Second)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}

func TestWaitWithSomeTimeAlreadyElapsed(t *testing.T) {
	// given
	action := game.Wait(3 * time.Second)
	action.Run(time.Second)
	statusWanted := game.Done(0)

	// when
	status := action.Run(2 * time.Second)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}

func TestProgress(t *testing.T) {
	tests := map[string]struct {
		Needed time.Duration
		Steps  []time.Duration

		Status   game.ActionStatus
		Progress float64
	}{
		"Start": {
			Needed: 3 * time.Second,
		},
		"Start/Instant": {
			Needed: 0,

			Progress: 1,
		},

		"Halfway": {
			Needed: 2 * time.Second,
			Steps: []time.Duration{
				time.Second,
			},

			Progress: 0.5,
		},
		"Halfway/InTwoSteps": {
			Needed: 2 * time.Second,
			Steps: []time.Duration{
				time.Second / 2,
				time.Second / 2,
			},

			Progress: 0.5,
		},

		"Complete": {
			Needed: time.Second,
			Steps: []time.Duration{
				time.Second,
			},

			Status:   game.Done(0),
			Progress: 1,
		},
		"Complete/InTwoSteps": {
			Needed: 2 * time.Second,
			Steps: []time.Duration{
				time.Second,
				time.Second,
			},

			Status:   game.Done(0),
			Progress: 1,
		},

		"AfterDone": {
			Needed: time.Second,
			Steps: []time.Duration{
				time.Second,
				time.Second,
			},

			Status:   game.Done(time.Second),
			Progress: 1,
		},
		"AfterDone/WithTimeToSpare": {
			Needed: time.Second,
			Steps: []time.Duration{
				2 * time.Second,
			},

			Status:   game.Done(time.Second),
			Progress: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			progress := game.Progress(tt.Needed)

			// when
			var status game.ActionStatus
			for _, dt := range tt.Steps {
				status = progress.Run(dt)
			}

			// then
			assert.That(
				status == tt.Status,
				t.Errorf, "got status %#v, want %#v", status, tt.Status)
			assert.That(
				progress.Progress() == tt.Progress,
				t.Errorf, "got progress %f, want %f", progress.Progress(), tt.Progress)
		})
	}
}

func TestSequence(t *testing.T) {
	tests := map[string]struct {
		Steps  []game.Action
		Times  []time.Duration
		Status game.ActionStatus
	}{
		"Empty": {
			Times:  []time.Duration{time.Second},
			Status: game.Done(time.Second),
		},

		"SingleStep/InsufficientTime": {
			Steps:  []game.Action{game.Wait(3 * time.Second)},
			Times:  []time.Duration{2 * time.Second},
			Status: game.Paused(),
		},
		"SingleStep/ExactTime": {
			Steps:  []game.Action{game.Wait(time.Second)},
			Times:  []time.Duration{time.Second},
			Status: game.Done(0),
		},
		"SingleStep/TimeToSpare": {
			Steps:  []game.Action{game.Wait(time.Second)},
			Times:  []time.Duration{3 * time.Second},
			Status: game.Done(2 * time.Second),
		},

		"TwoSteps/First/InsufficientTime": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{time.Second},
			Status: game.Paused(),
		},
		"TwoSteps/First/ExactTime": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{2 * time.Second},
			Status: game.Paused(),
		},
		"TwoSteps/Second/InsufficientTime": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{4 * time.Second},
			Status: game.Paused(),
		},
		"TwoSteps/Second/ExactTime": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{5 * time.Second},
			Status: game.Done(0),
		},
		"TwoSteps/Second/TimeToSpare": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{6 * time.Second},
			Status: game.Done(time.Second),
		},

		"Interrupted": {
			Steps: []game.Action{
				game.Interrupt(),
				game.Wait(time.Second),
			},
			Times:  []time.Duration{2 * time.Second},
			Status: game.Interrupted(2 * time.Second),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			assert.That(len(tt.Times) > 0, t.Errorf, "case %q has no times specified", name)
			action := game.Sequence(tt.Steps...)

			// when
			var status game.ActionStatus
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
