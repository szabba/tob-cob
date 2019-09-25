// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui_test

import (
	"testing"

	"github.com/faiface/pixel"
	"github.com/szabba/assert"

	"github.com/szabba/tob-cob/ui"
)

func TestZeroCameraCentersTheOrigin(t *testing.T) {
	// given
	bounds := pixel.R(0, 0, 800, 600)
	center := bounds.Center()

	// when
	cam := ui.Camera{}

	// then
	originAt := cam.Matrix(bounds).Project(pixel.ZV)
	assert.That(
		originAt == center,
		t.Errorf, "origin at %s, want it at %s", originAt, center)
}

func TestCameraCentersTheLookAtPoint(t *testing.T) {
	// given
	bounds := pixel.R(0, 0, 800, 600)
	center := bounds.Center()
	lookAt := pixel.V(300, 200)

	// when
	cam := ui.NewCamera(lookAt)

	// then
	onscreen := cam.Matrix(bounds).Project(lookAt)
	assert.That(
		onscreen == center,
		t.Errorf, "origin at %s, want it at %s", onscreen, center)
}

func TestCameraMoveByShiftsThePointBeingLookedAt(t *testing.T) {
	// given
	bounds := pixel.R(0, 0, 800, 600)
	center := bounds.Center()
	lookAt := pixel.V(300, 200)
	shift := pixel.V(100, 50)
	cam := ui.NewCamera(lookAt)

	// when
	cam.MoveBy(shift)

	// then
	onscreen := cam.Matrix(bounds).Project(lookAt.Add(shift))
	assert.That(
		onscreen == center,
		t.Errorf, "origin at %s, want it at %s", onscreen, center)
}
