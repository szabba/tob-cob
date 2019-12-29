// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game_test

import (
	"testing"

	"github.com/szabba/assert"
	"github.com/szabba/tob-cob/game"
)

func TestOnePosTakerStartsEmpty(t *testing.T) {
	// given
	// when
	taker := game.OnePosTaker{}

	// then
	assert.That(!taker.Placed(), t.Errorf, "the taker is placed - it should not")
	assert.That(
		taker.AtPoint() == game.Point{},
		t.Errorf, "reported at point %#v - want %#v", taker.AtPoint(), game.Point{})
}

func TestOnePosTakerCanTakePosition(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Create()

	taker := game.OnePosTaker{}

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
	space := game.NewSpace()
	src := space.At(game.P(2, 3))
	src.Create()
	dst := space.At(game.P(2, 4))
	dst.Create()

	taker := game.OnePosTaker{}
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
	space := game.NewSpace()
	src := space.At(game.P(2, 3))
	src.Create()
	dst := space.At(game.P(2, 4))
	dst.Create()

	taker := game.OnePosTaker{}
	src.Take(&taker)

	// when
	dst.Take(&taker)

	// then
	assert.That(!src.Taken(), t.Errorf, "the original position is still taken")
}

func TestOnePosTakerCanLeaveTakenPosition(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Create()

	taker := game.OnePosTaker{}
	pos.Take(&taker)

	// when
	taker.Leave()

	// then
	assert.That(!taker.Placed(), t.Errorf, "the taker is placed - it should not be")
	assert.That(
		taker.AtPoint() == game.Point{},
		t.Fatalf, "reported at point %#v - want %#v", taker.AtPoint(), game.Point{})
}
