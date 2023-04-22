// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package grid_test

import (
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/tob-cob/game/grid"
)

func TestOnePosTakerStartsEmpty(t *testing.T) {
	// given
	// when
	taker := grid.OnePosTaker{}

	// then
	assert.That(!taker.Placed(), t.Errorf, "the taker is placed - it should not")
	assert.That(
		taker.AtPoint() == grid.Point{},
		t.Errorf, "reported at point %#v - want %#v", taker.AtPoint(), grid.Point{})
}

func TestOnePosTakerCanTakePosition(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()

	taker := grid.OnePosTaker{}

	// when
	pos.Take(&taker)

	// then
	assert.That(taker.Placed(), t.Errorf, "the taker is not placed - it should")
	assert.That(
		taker.AtPoint() == pos.AtPoint(),
		t.Fatalf, "reported at point %#v - want %#v", taker.AtPoint(), pos.AtPoint())
}

func TestOnePosTakerCanBeMovedToNewPosition(t *testing.T) {
	// given
	space := grid.NewSpace()
	src := space.At(grid.P(2, 3))
	src.Create()
	dst := space.At(grid.P(2, 4))
	dst.Create()

	taker := grid.OnePosTaker{}
	src.Take(&taker)

	// when
	dst.Take(&taker)

	// then
	assert.That(taker.Placed(), t.Errorf, "the taker is not placed - it should")
	assert.That(
		taker.AtPoint() == dst.AtPoint(),
		t.Fatalf, "reported at point %#v - want %#v", taker.AtPoint(), dst.AtPoint())
}

func TestOnePosTakerLeavesOldPosition(t *testing.T) {
	// given
	space := grid.NewSpace()
	src := space.At(grid.P(2, 3))
	src.Create()
	dst := space.At(grid.P(2, 4))
	dst.Create()

	taker := grid.OnePosTaker{}
	src.Take(&taker)

	// when
	dst.Take(&taker)

	// then
	assert.That(!src.Taken(), t.Errorf, "the original position is still taken")
}

func TestOnePosTakerCanLeaveTakenPosition(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()

	taker := grid.OnePosTaker{}
	pos.Take(&taker)

	// when
	taker.Leave()

	// then
	assert.That(!taker.Placed(), t.Errorf, "the taker is placed - it should not be")
	assert.That(
		taker.AtPoint() == grid.Point{},
		t.Fatalf, "reported at point %#v - want %#v", taker.AtPoint(), grid.Point{})
}
