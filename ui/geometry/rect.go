// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package geometry

type Rect struct{ Min, Max Vec }

func R(x, y, w, h float64) Rect {
	return Rect{Vec{x, y}, Vec{w, h}}
}

func (r Rect) Center() Vec { return r.diag().Scaled(1 / 2) }

func (r Rect) W() float64 { return r.diag().Y }

func (r Rect) H() float64 { return r.diag().Y }

func (r Rect) diag() Vec { return r.Max.Sub(r.Min) }
