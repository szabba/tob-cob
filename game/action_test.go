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
		!status.HasTimeLeft(),
		t.Errorf, "the status should say there is no time left after running the action")
	assert.That(
		status.TimeLeft() == timeLeft,
		t.Errorf, "time left is %s, want %s", status.TimeLeft(), timeLeft)
}

func TestDoneWithNegaiveTimeLeft(t *testing.T) {
	// given
	timeLeft := -time.Second
	sanitizedTime := time.Duration(0)

	// when
	status := game.Done(timeLeft)

	// then
	assert.That(status.Done(), t.Errorf, "the status should say the action has completed")
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

type TestAction struct {
	OnError         func(msg string, args ...interface{})
	RunTimes        int
	TimeStillNeeded time.Duration
}

var _ game.Action = &TestAction{}

func (action *TestAction) Run(timeLeft time.Duration) game.ActionStatus {
	action.RunTimes++
	if timeLeft < 0 && action.OnError != nil {
		action.OnError(
			"action %p (in state %#v) called with timeLeft %s < 0",
			action, action, timeLeft)
	}

	if timeLeft < action.TimeStillNeeded {
		action.TimeStillNeeded -= timeLeft
		timeLeft = 0
		return game.Paused()
	} else {
		action.TimeStillNeeded = 0
		timeLeft -= action.TimeStillNeeded
		return game.Done(timeLeft)
	}
}
