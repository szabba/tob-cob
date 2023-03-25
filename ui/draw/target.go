// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package draw

import (
	"image"
	"image/color"

	"github.com/szabba/tob-cob/ui/geometry"
)

type Target interface {
	Clear(c color.Color)
	SetMatrix(m geometry.Mat)
	Rectangle(r geometry.Rect, m geometry.Mat, fill color.Color)
	Import(img image.Image) Image
}

type Image interface {
	Bounds() geometry.Rect
	Draw(m geometry.Mat)
}
