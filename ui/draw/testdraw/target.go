// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testdraw

import (
	"image"
	"image/color"

	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/geometry"
)

type Target struct{}

var _ draw.Target = new(Target)

func (*Target) Clear(_ color.Color)      {}
func (*Target) SetMatrix(_ geometry.Mat) {}

func (*Target) Import(img image.Image) draw.Image {
	return _Image{
		img: img,
	}
}

type _Image struct {
	img image.Image
}

var _ draw.Image = new(_Image)

func (img _Image) Bounds() geometry.Rect {
	stdlib := img.img.Bounds()
	return geometry.R(
		0, 0,
		float64(stdlib.Dx()),
		float64(stdlib.Dy()))
}

func (_Image) Draw(_ geometry.Mat, _ geometry.Vec) {}
