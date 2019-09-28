// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/rs/zerolog/log"
)

// A Camera describes the coordinate tranformation between the world and the window.
//
// The zero value looks at the origin point (pixel.ZV) in world coordinates.
type Camera struct {
	lookAt pixel.Vec
}

// NewCamera creates a camera that will put lookAt in the center of the window.
func NewCamera(lookAt pixel.Vec) Camera {
	return Camera{lookAt}
}

// MoveBy changes the point being looked at by delta in window coordinates.
func (cam *Camera) MoveBy(delta pixel.Vec) {
	cam.lookAt = cam.lookAt.Add(delta)
}

// Matrix computes the world-to-window coordinate transformation matrix.
func (cam *Camera) Matrix(bounds pixel.Rect) pixel.Matrix {
	center := bounds.Center()
	matrix := pixel.IM.Moved(center).Moved(cam.lookAt.Scaled(-1))
	log.Debug().
		Float64("center.x", center.X).
		Float64("center.y", center.Y).
		Float64("lookAt.x", cam.lookAt.X).
		Float64("lookAt.y", cam.lookAt.Y).
		Str("matrix", matrix.String()).
		Msg("camera matrix calculated")
	return matrix
}

type CameraController struct {
	cam *Camera
}

func NewCamController(cam *Camera) *CameraController {
	return &CameraController{cam: cam}
}

func (cont *CameraController) Process(input Input) {
	delta := cont.lookAtDelta(input)
	cont.cam.MoveBy(delta)
}

func (cont *CameraController) lookAtDelta(input Input) pixel.Vec {
	delta := pixel.ZV
	if input.Pressed(pixelgl.KeyLeft) {
		delta.X -= 5
	}
	if input.Pressed(pixelgl.KeyRight) {
		delta.X += 5
	}
	if input.Pressed(pixelgl.KeyUp) {
		delta.Y += 5
	}
	if input.Pressed(pixelgl.KeyDown) {
		delta.Y -= 5
	}
	return delta
}
