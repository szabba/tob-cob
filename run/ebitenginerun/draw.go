// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ebitenginerun

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/geometry"
)

var (
	dst = _DrawTarget{}

	_ draw.Target = &_DrawTarget{}
)

type _DrawTarget struct {
	dst    *ebiten.Image
	bounds geometry.Rect
	matrix geometry.Mat
}

func (d *_DrawTarget) _SetBounds(width, height int) {
	d.bounds = geometry.R(0, 0, float64(width), float64(height))
}

func (d *_DrawTarget) Clear(c color.Color) { d.dst.Fill(c) }

func (d *_DrawTarget) SetMatrix(m geometry.Mat) { d.matrix = m }

func (d *_DrawTarget) Import(img image.Image) draw.Image {
	src := ebiten.NewImageFromImage(img)
	return _DrawImage{d, src}
}

type _DrawImage struct {
	dst *_DrawTarget
	src *ebiten.Image
}

func (d _DrawImage) Bounds() geometry.Rect {

	bounds := d.src.Bounds()

	return geometry.R(
		float64(bounds.Min.X),
		float64(bounds.Min.Y),
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}

func (d _DrawImage) Draw(m geometry.Mat) {
	opts := &ebiten.DrawImageOptions{
		GeoM: d.toGeoM(d.dst.matrix.Compose(m)),
	}
	d.dst.dst.DrawImage(d.src, opts)
}

func (d _DrawImage) toGeoM(m geometry.Mat) ebiten.GeoM {

	composed := d.postFixM().Compose(m).Compose(d.preFixM())

	return d.toGeoMLiteral(composed)
}

func (_DrawImage) toGeoMLiteral(m geometry.Mat) ebiten.GeoM {
	var out ebiten.GeoM
	for i, row := range m {
		for j, m_ij := range row {
			out.SetElement(i, j, m_ij)
		}
	}
	return out
}

func (d _DrawImage) preFixM() geometry.Mat {
	bounds := d.Bounds()
	return geometry.Mat{
		{1, 0, 0},
		{0, -1, bounds.H()},
	}
}

func (d _DrawImage) postFixM() geometry.Mat {
	h := d.dst.bounds.H()
	return geometry.Mat{
		{1, 0, 0},
		{0, -1, h},
	}
}
