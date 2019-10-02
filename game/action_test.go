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
