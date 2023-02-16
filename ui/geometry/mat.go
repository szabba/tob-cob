// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package geometry

import "fmt"

type Mat [2][3]float64

func (m Mat) Zero() bool { return m == Mat{} }

func (m Mat) String() string { return fmt.Sprintf("%v", [2][3]float64(m)) }

func Translation(offset Vec) Mat {
	t := Identity()
	t[0][2] = offset.X
	t[1][2] = offset.Y
	return t
}

func Identity() Mat { return Scale(1) }

func Scale(a float64) Mat {
	s := Mat{}
	s[0][0] = a
	s[1][1] = a
	return s
}

func (m Mat) Compose(o Mat) Mat {
	out := Mat{}
	// TODO: Finish
	return out
}

func (m Mat) Apply(v Vec) Vec {
	return Vec{
		X: m.applyRow(0, v),
		Y: m.applyRow(1, v),
	}
}

func (m Mat) applyRow(i int, v Vec) float64 {
	return m[i][0]*v.X + m[i][1]*v.Y + m[i][2]
}

func (m Mat) Invert() Mat {
	// I had better things to do than derive this...
	// https://www.wolframalpha.com/input?i=Inverse%5B%7B%7Ba%2C+b%2C+c%7D%2C+%7Bd%2C+e%2C+f%7D%2C+%7B0%2C+0%2C+1%7D%7D%5D

	a, b, c := m[0][0], m[0][1], m[0][2]
	d, e, f := m[1][0], m[1][1], m[1][2]

	return Mat{
		{
			e / (a*e - b*d),
			b / (b*d - a*e),
			(c*e - b*f) / (b*d - a*e),
		},
		{
			d / (b*d - a*e),
			a / (a*e - b*d),
			(c*d - a*f) / (a*e - b*d),
		},
	}
}
