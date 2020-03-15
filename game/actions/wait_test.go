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

func TestWaitWithExactTime(t *testing.T) {
	// given
	action := actions.Wait(time.Second)
	statusWanted := actions.Done(0)

	// when
	status := action.Run(time.Second)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}

func TestWaitWithoutEnoughTime(t *testing.T) {
	// given
	action := actions.Wait(3 * time.Second)
	statusWanted := actions.Paused()

	// when
	status := action.Run(time.Second)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}
func TestWaitWithTimeToSpare(t *testing.T) {
	// given
	action := actions.Wait(time.Second)
	statusWanted := actions.Done(2 * time.Second)

	// when
	status := action.Run(3 * time.Second)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}

func TestWaitWithSomeTimeAlreadyElapsed(t *testing.T) {
	// given
	action := actions.Wait(3 * time.Second)
	action.Run(time.Second)
	statusWanted := actions.Done(0)

	// when
	status := action.Run(2 * time.Second)

	// then
	assert.That(
		status == statusWanted,
		t.Errorf, "got status %#v, want %#v", status, statusWanted)
}
