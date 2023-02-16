// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"image"
	"math"
)

type Rect struct{ Min, Max Vec }

func RectFromImageRect(r image.Rectangle) Rect {
	// TODO: FIXME?
	return Rect{
		Min: Vec{
			X: float64(r.Min.X),
			Y: float64(r.Bounds().Dy() - r.Min.Y),
		},
		Max: Vec{
			X: float64(r.Max.X),
			Y: float64(r.Bounds().Dy() - r.Max.Y),
		},
	}
}

func (r Rect) Moved(v Vec) Rect {
	r.Min = r.Min.Add(v)
	r.Max = r.Max.Add(v)
	return r
}

func (r Rect) W() float64 { return math.Abs(r.Max.X - r.Min.X) }
func (r Rect) H() float64 { return math.Abs(r.Max.Y - r.Min.Y) }

func (r Rect) Center() Vec {
	return Vec{
		(r.Min.X + r.Max.X) / 2,
		(r.Min.Y + r.Min.Y) / 2,
	}
}

type Vec struct{ X, Y float64 }

func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y}
}

func (v Vec) Scaled(a float64) Vec {
	return Vec{a * v.X, a * v.Y}
}
