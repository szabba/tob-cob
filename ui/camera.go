// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
	"golang.org/x/exp/slog"
)

// A Camera describes the coordinate tranformation between the world and the window.
//
// The zero value looks at the origin point (pixel.ZV) in world coordinates.
type Camera struct {
	lookAt geometry.Vec
}

// NewCamera creates a camera that will put lookAt in the center of the window.
func NewCamera(lookAt geometry.Vec) Camera {
	return Camera{lookAt}
}

// MoveBy changes the point being looked at by delta in window coordinates.
func (cam *Camera) MoveBy(delta geometry.Vec) {
	cam.lookAt = cam.lookAt.Add(delta)

}

// Matrix computes the world-to-window coordinate transformation matrix.
func (cam *Camera) Matrix(bounds geometry.Rect) geometry.Mat {
	center := bounds.Center()
	matrix := geometry.Translation(center.Add(cam.lookAt.Scaled(-1)))

	if slog.Default().Enabled(nil, slog.LevelDebug) {
		slog.Debug(
			"camera matrix calculated",

			slog.Group("center",
				slog.Float64("x", center.X),
				slog.Float64("y", center.Y),
			),

			slog.Group("look-at",
				slog.Float64("x", cam.lookAt.X),
				slog.Float64("y", cam.lookAt.Y),
			),

			slog.String("matrix", matrix.String()),
		)
	}

	return matrix
}

type CameraController struct {
	cam *Camera
}

func NewCamController(cam *Camera) *CameraController {
	return &CameraController{cam: cam}
}

func (cont *CameraController) Process(src input.Source) {
	if !src.Focused() {
		return
	}
	delta := cont.lookAtDelta(src)
	cont.cam.MoveBy(delta)
}

func (cont *CameraController) lookAtDelta(src input.Source) geometry.Vec {
	if !src.MouseInsideWindow() {
		return geometry.Vec{}
	}

	delta := geometry.Vec{}
	if src.Pressed(input.KeyLeft()) || cont.mouseNearLeftEdge(src) {
		delta.X -= 5
	}
	if src.Pressed(input.KeyRight()) || cont.mouseNearRightEdge(src) {
		delta.X += 5
	}
	if src.Pressed(input.KeyUp()) || cont.mouseNearTopEdge(src) {
		delta.Y += 5
	}
	if src.Pressed(input.KeyDown()) || cont.mouseNearBottomEdge(src) {
		delta.Y -= 5
	}
	return delta
}

func (*CameraController) mouseNearLeftEdge(src input.Source) bool {
	width := src.Bounds().W()
	return src.MousePosition().X < 0.05*width
}

func (*CameraController) mouseNearRightEdge(src input.Source) bool {
	width := src.Bounds().W()
	return src.MousePosition().X > 0.95*width
}

func (*CameraController) mouseNearBottomEdge(src input.Source) bool {
	height := src.Bounds().H()
	return src.MousePosition().Y < 0.05*height
}

func (*CameraController) mouseNearTopEdge(src input.Source) bool {
	height := src.Bounds().H()
	return src.MousePosition().Y > 0.95*height
}
