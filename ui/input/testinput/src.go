// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testinput

import (
	"github.com/szabba/tob-cob/ui/geometry"
	"github.com/szabba/tob-cob/ui/input"
)

type Source struct {
	Mock struct {
		Bounds            func() geometry.Rect
		MousePosition     func() geometry.Vec
		MouseInsideWindow func() bool
	}
}

var _ input.Source = Source{}

func (Source) Focused() bool { return false }

func (Source) JustReleased(_ input.Button) bool { return false }

func (Source) JustPressed(_ input.Button) bool { return false }

func (Source) Pressed(_ input.Button) bool { return false }

func (src Source) Bounds() geometry.Rect {
	if src.Mock.Bounds == nil {
		return geometry.Rect{}
	}
	return src.Mock.Bounds()
}

func (src Source) MouseInsideWindow() bool {
	if src.Mock.MouseInsideWindow == nil {
		return false
	}
	return src.Mock.MouseInsideWindow()
}

func (src Source) MousePosition() geometry.Vec {
	if src.Mock.MousePosition == nil {
		return geometry.Vec{}
	}
	return src.Mock.MousePosition()
}
