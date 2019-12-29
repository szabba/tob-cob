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

func TestCountdown(t *testing.T) {
	tests := map[string]struct {
		Needed int
		Steps  []int

		LeftOver int
		Progress float64
	}{
		"Start": {
			Needed: 3,
		},
		"Start/Instant": {
			Needed: 0,

			Progress: 1,
		},

		"Halfway": {
			Needed: 2,
			Steps:  []int{1},

			Progress: 0.5,
		},
		"Halfway/InTwoSteps": {
			Needed: 4,
			Steps:  []int{1, 1},

			Progress: 0.5,
		},

		"Complete": {
			Needed: 1,
			Steps:  []int{1},

			Progress: 1,
		},
		"Complete/InTwoSteps": {
			Needed: 2,
			Steps:  []int{1, 1},

			Progress: 1,
		},

		"AfterDone": {
			Needed: 1,
			Steps:  []int{1, 1},

			LeftOver: 1,
			Progress: 1,
		},
		"AfterDone/WithTimeToSpare": {
			Needed: 1,
			Steps:  []int{2},

			LeftOver: 1,
			Progress: 1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			countdown := game.CountdownTo(tt.Needed)

			// when
			var leftOver int
			for _, step := range tt.Steps {
				leftOver = countdown.CountDown(step)
			}

			// then
			assert.That(
				leftOver == tt.LeftOver,
				t.Errorf, "got left over steps %#v, want %#v", leftOver, tt.LeftOver)
			assert.That(
				countdown.Progress() == tt.Progress,
				t.Errorf, "got progress %f, want %f", countdown.Progress(), tt.Progress)
		})
	}
}

func TestCountdownOver(t *testing.T) {
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
			countdown, action := game.CountdownOver(tt.Needed)

			// when
			var status game.ActionStatus
			for _, dt := range tt.Steps {
				status = action.Run(dt)
			}

			// then
			assert.That(
				status == tt.Status,
				t.Errorf, "got status %#v, want %#v", status, tt.Status)
			assert.That(
				countdown.Progress() == tt.Progress,
				t.Errorf, "got progress %f, want %f", countdown.Progress(), tt.Progress)
		})
	}
}