// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui_test

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/szabba/tob-cob/ui"
)

type TestInput struct {
	Mock struct {
		Bounds            func() pixel.Rect
		MousePosition     func() pixel.Vec
		MouseInsideWindow func() bool
	}
}

var _ ui.Input = TestInput{}

func (input TestInput) Focused() bool { return false }

func (input TestInput) Pressed(pixelgl.Button) bool { return false }

func (input TestInput) Bounds() pixel.Rect {
	if input.Mock.Bounds == nil {
		return pixel.R(0, 0, 0, 0)
	}
	return input.Mock.Bounds()
}

func (input TestInput) MouseInsideWindow() bool {
	if input.Mock.MouseInsideWindow == nil {
		return false
	}
	return input.Mock.MouseInsideWindow()
}

func (input TestInput) MousePosition() pixel.Vec {
	if input.Mock.MousePosition == nil {
		return pixel.V(0, 0)
	}
	return input.Mock.MousePosition()
}
