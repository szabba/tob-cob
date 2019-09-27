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
	cam          *Camera
	processInput inputProcessor
}

func NewCamController(cam *Camera) *CameraController {
	return &CameraController{cam: cam}
}

type inputProcessor func(Input) inputProcessor

func (cont *CameraController) Process(input Input) {
	if cont.processInput == nil {
		cont.processInput = cont.start(input)
	} else {
		cont.processInput = cont.processInput(input)
	}
}

func (cont *CameraController) start(input Input) inputProcessor {
	if !input.Pressed(pixelgl.MouseButtonLeft) {
		return nil
	}
	return cont.moving(input)
}

func (cont *CameraController) moving(input Input) inputProcessor {
	from := input.MousePosition()
	return func(input Input) inputProcessor {
		if !input.Pressed(pixelgl.MouseButtonLeft) {
			return nil
		}

		delta := input.MousePosition().Sub(from).Scaled(-1)
		cont.cam.MoveBy(delta)
		return cont.moving(input)
	}
}
