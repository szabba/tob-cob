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
	Import(img image.Image) Image
}

type Image interface {
	Bounds() geometry.Rect

	// Draw the image on the target that imported it.
	// The drawn image is transforming by the matrix m and then by the current target-wide matrix.
	// The anchor is understood to use the same coordinate system as the bounds do.
	Draw(m geometry.Mat, anchor geometry.Vec)
}
