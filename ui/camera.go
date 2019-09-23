// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"github.com/faiface/pixel"
	"github.com/rs/zerolog/log"
)

type Camera struct {
	lookAt pixel.Vec
}

func NewCamera(lookAt pixel.Vec) Camera {
	return Camera{lookAt}
}

func (cam *Camera) MoveBy(delta pixel.Vec) {
	cam.lookAt = cam.lookAt.Add(delta)
}

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
