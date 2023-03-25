// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package input

import (
	"github.com/szabba/tob-cob/ui/geometry"
)

type Source interface {
	Focused() bool
	JustReleased(btn Button) bool
	JustPressed(btn Button) bool
	Pressed(btn Button) bool
	MousePosition() geometry.Vec
	MouseInsideWindow() bool
	Bounds() geometry.Rect
}
