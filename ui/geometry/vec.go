// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package geometry

type Vec struct{ X, Y float64 }

func V(x, y float64) Vec { return Vec{x, y} }

func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y}
}

func (v Vec) Sub(o Vec) Vec {
	return Vec{v.X - o.X, v.Y - o.Y}
}

func (v Vec) Scaled(s float64) Vec {
	return Vec{s * v.X, s * v.Y}
}
