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

	// when
	status := action.Run(timeLeft)

	// then
	assert.That(status.Done(), t.Errorf, "the action should have completed")
	assert.That(
		status.TimeLeft() == timeLeft,
		t.Errorf, "the action left %s, wanted %s", status.TimeLeft(), timeLeft)
}

func TestWaitWithExactTime(t *testing.T) {
	// given
	action := game.Wait(time.Second)

	// when
	status := action.Run(time.Second)

	// then
	assert.That(status.Done(), t.Errorf, "the action should complete")
	assert.That(!status.HasTimeLeft(), t.Errorf, "the action should use up all time")
	assert.That(
		status.TimeLeft() == time.Duration(0),
		t.Errorf, "the time left is %s, want %s", status.TimeLeft(), time.Duration(0))
}

func TestWaitWithoutEnoughTime(t *testing.T) {
	// given
	action := game.Wait(3 * time.Second)

	// when
	status := action.Run(time.Second)

	// then
	assert.That(!status.Done(), t.Errorf, "the action should not complete")
	assert.That(!status.HasTimeLeft(), t.Errorf, "the action should use up all time")
	assert.That(
		status.TimeLeft() == time.Duration(0),
		t.Errorf, "the time left is %s, want %s", status.TimeLeft(), time.Duration(0))
}
func TestWaitWithTimeToSpare(t *testing.T) {
	// given
	action := game.Wait(time.Second)

	// when
	status := action.Run(3 * time.Second)

	// then
	assert.That(status.Done(), t.Errorf, "the action should complete")
	assert.That(status.HasTimeLeft(), t.Errorf, "the action should leave some time")
	assert.That(
		status.TimeLeft() == 2*time.Second,
		t.Errorf, "the time left is %s, want %s", status.TimeLeft(), 2*time.Second)
}

func TestWaitWithSomeTimeAlreadyElapsed(t *testing.T) {
	// given
	action := game.Wait(3 * time.Second)
	action.Run(time.Second)

	// when
	status := action.Run(2 * time.Second)

	// then
	assert.That(status.Done(), t.Errorf, "the action should complete")
	assert.That(!status.HasTimeLeft(), t.Errorf, "the action should leave no time")
	assert.That(
		status.TimeLeft() == time.Duration(0),
		t.Errorf, "the time left is %s, want %s", status.TimeLeft(), time.Duration(0))
}

func TestSequence(t *testing.T) {
	kases := map[string]struct {
		Steps  []game.Action
		Times  []time.Duration
		Status game.ActionStatus
	}{
		"empty": {
			Times:  []time.Duration{time.Second},
			Status: game.Done(time.Second),
		},

		"singleStep/insufficientTime": {
			Steps:  []game.Action{game.Wait(3 * time.Second)},
			Times:  []time.Duration{2 * time.Second},
			Status: game.Paused(),
		},
		"singleStep/exactTime": {
			Steps:  []game.Action{game.Wait(time.Second)},
			Times:  []time.Duration{time.Second},
			Status: game.Done(0),
		},
		"singleStep/timeToSpare": {
			Steps:  []game.Action{game.Wait(time.Second)},
			Times:  []time.Duration{3 * time.Second},
			Status: game.Done(2 * time.Second),
		},

		"twoSteps/first/insufficientTime": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{time.Second},
			Status: game.Paused(),
		},
		"twoSteps/first/exactTime": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{2 * time.Second},
			Status: game.Paused(),
		},
		"twoSteps/second/insufficientTime": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{4 * time.Second},
			Status: game.Paused(),
		},
		"twoSteps/second/exactTime": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{5 * time.Second},
			Status: game.Done(0),
		},
		"twoSteps/second/timeToSpare": {
			Steps: []game.Action{
				game.Wait(2 * time.Second),
				game.Wait(3 * time.Second),
			},
			Times:  []time.Duration{6 * time.Second},
			Status: game.Done(time.Second),
		},
	}

	for name, kase := range kases {
		t.Run(name, func(t *testing.T) {
			// given
			assert.That(len(kase.Times) > 0, t.Errorf, "case %q has no times specified", name)
			action := game.Sequence(kase.Steps...)
			var status game.ActionStatus

			// when
			for _, t := range kase.Times {
				status = action.Run(t)
			}

			// then
			assert.That(
				status == kase.Status,
				t.Errorf, "final action status is %#v, want %#v", status, kase.Status)
		})
	}
}
