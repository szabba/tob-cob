// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Input interface {
	Focused() bool
	Pressed(btn pixelgl.Button) bool
	MousePosition() pixel.Vec
	MouseInsideWindow() bool
	Bounds() pixel.Rect
}

var _ Input = &pixelgl.Window{}
