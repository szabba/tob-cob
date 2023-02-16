// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui_test

import (
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/tob-cob/ui"
	"github.com/szabba/tob-cob/ui/geometry"
)

func TestZeroCameraCentersTheOrigin(t *testing.T) {
	// given
	bounds := geometry.R(0, 0, 800, 600)
	center := bounds.Center()

	// when
	cam := ui.Camera{}

	// then
	originAt := cam.Matrix(bounds).Apply(geometry.Vec{})
	assert.That(
		originAt == center,
		t.Errorf, "origin at %s, want it at %s", originAt, center)
}

func TestCameraCentersTheLookAtPoint(t *testing.T) {
	// given
	bounds := geometry.R(0, 0, 800, 600)
	center := bounds.Center()
	lookAt := geometry.V(300, 200)

	// when
	cam := ui.NewCamera(lookAt)

	// then
	onscreen := cam.Matrix(bounds).Apply(lookAt)
	assert.That(
		onscreen == center,
		t.Errorf, "origin at %s, want it at %s", onscreen, center)
}

func TestCameraMoveByShiftsThePointBeingLookedAt(t *testing.T) {
	// given
	bounds := geometry.R(0, 0, 800, 600)
	center := bounds.Center()
	lookAt := geometry.V(300, 200)
	shift := geometry.V(100, 50)
	cam := ui.NewCamera(lookAt)

	// when
	cam.MoveBy(shift)

	// then
	onscreen := cam.Matrix(bounds).Apply(lookAt.Add(shift))
	assert.That(
		onscreen == center,
		t.Errorf, "origin at %s, want it at %s", onscreen, center)
}
