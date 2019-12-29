// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game

// A OnePosTaker is a space taker that occupies at most one grid position at a time.
type OnePosTaker struct {
	pos Position
}

// LetOnto is part of the SpaceTaker interface.
func (taker *OnePosTaker) LetOnto(pos Position) {
	if taker.Placed() {
		taker.pos.Free()
	}
	taker.pos = pos
}

// ForceOff is part of the SpaceTaker interface.
func (taker *OnePosTaker) ForceOff(pos Position) {
}

// Placed says whether the space taker is taking a position.
func (taker *OnePosTaker) Placed() bool {
	return taker.pos != Position{}
}

// AtPoint is the point at which the taken position is.
//
// The zero value of the type is returned when the taker is not placed.
// Do not use that to check if the taker is placed.
// A taker might be taking the position at the zero point.
func (taker *OnePosTaker) AtPoint() Point { return taker.pos.AtPoint() }

// Leave makes the taker leave a position if it has one taken.
func (taker *OnePosTaker) Leave() {
	if !taker.Placed() {
		return
	}
	taker.pos.Free()
	taker.pos = Position{}
}
