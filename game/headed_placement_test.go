// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game_test

import (
	"testing"
	"time"

	"github.com/szabba/assert"
	"github.com/szabba/tob-cob/game"
)

func TestZeroHeadedPlacementIsNotPlaced(t *testing.T) {
	// given
	// when
	zero := game.HeadedPlacement{}

	// then
	assertNotPlaced(t, &zero)
}

func TestZeroHeadedPlacementIsNotHeaded(t *testing.T) {
	// given
	// when
	zero := game.HeadedPlacement{}

	// then
	assertNotHeaded(t, &zero)
}

func TestZeroHeadedPlacementHasNoProgress(t *testing.T) {
	// given
	// when
	zero := game.HeadedPlacement{}

	// then
	assert.That(zero.Progress() == 0, t.Errorf, "got progress %v - want %v", zero.Progress(), 0)
}

func TestHeadedPlacementCannotBePlacedAtPositionThatDoesNotExist(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))

	placement := game.HeadedPlacement{}

	// when
	ok := placement.Place(pos)

	// then
	assert.That(!ok, t.Errorf, "placing should fail")
	assertNotPlaced(t, &placement)
	assert.That(!pos.Taken(), t.Errorf, "the position is taken - it should not be")
}

func TestHeadedPlacementCannotBePlacedAtPositionThatIsAlreadyTaken(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Take(game.DummyTaker())

	placement := game.HeadedPlacement{}

	// when
	ok := placement.Place(pos)

	// then
	assert.That(!ok, t.Errorf, "placing should fail")
	assertNotPlaced(t, &placement)
	assert.That(!pos.Taken(), t.Errorf, "the position is taken - it should not be")
}

func TestHeadedPlacementCanBePlaced(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Create()

	placement := game.HeadedPlacement{}

	// when
	ok := placement.Place(pos)

	// then
	assert.That(ok, t.Errorf, "placing should succeed")
	assertPlaced(t, &placement, pos.AtPoint())
	assert.That(pos.Taken(), t.Errorf, "the position is not taken - it should be")
	assert.That(
		placement.Progress() == 1,
		t.Errorf, "got progress %v, want %v", placement.Progress(), 1)
}

func TestHeadedPlacementCannotMoveIfNotPlaced(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Create()

	placement := game.HeadedPlacement{}

	// when
	action := placement.MoveTo(pos, time.Second)

	// then
	assert.That(
		action == game.NoAction(),
		t.Errorf, "got action %#v, want %#v", action, game.NoAction())
}

func TestPlacedHeadedPlacementCanMove(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Create()
	dst := space.At(game.P(2, 4))
	dst.Create()

	placement := game.HeadedPlacement{}
	placement.Place(pos)

	timeNeeded := time.Second
	action := placement.MoveTo(dst, timeNeeded)

	// when
	status := action.Run(timeNeeded)

	// then
	assert.That(
		status == game.Done(0),
		t.Errorf, "got move status %#v - want %#v", status, game.Done(0))
	assertPlaced(t, &placement, dst.AtPoint())
	assert.That(dst.Taken(), t.Errorf, "%#v should be taken, but is not", dst.AtPoint())
	assert.That(!pos.Taken(), t.Errorf, "%#v should not be taken, but it is", pos.AtPoint())
}

func TestPlacedHeadedPlacementBecomesHeadedWhenAMoveStarts(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Create()
	dst := space.At(game.P(2, 4))
	dst.Create()

	placement := game.HeadedPlacement{}
	placement.Place(pos)

	timeNeeded := time.Second
	action := placement.MoveTo(dst, timeNeeded)

	// when
	action.Run(0)

	// then
	assertHeaded(t, &placement, dst.AtPoint())
}

func TestHeadedPlacementProgressReflectsFractionOfTheMoveTimeThatHasPassed(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Create()
	dst := space.At(game.P(2, 4))
	dst.Create()

	placement := game.HeadedPlacement{}
	placement.Place(pos)

	timeNeeded := 4 * time.Second
	action := placement.MoveTo(dst, timeNeeded)

	// when
	action.Run(timeNeeded * 3 / 4)

	// then
	assert.That(placement.Progress() == 0.75, t.Errorf, "got progress %f - want %f", placement.Progress(), 0.75)
}

func TestFreshlyPlacedHeadedPlacementHasNoHeading(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(2, 3))
	pos.Create()

	placement := game.HeadedPlacement{}

	// when
	placement.Place(pos)

	// then
	assertPlaced(t, &placement, pos.AtPoint())
	assertNotHeaded(t, &placement)
}

func assertNotPlaced(t *testing.T, pl *game.HeadedPlacement) {
	assert.That(!pl.Placed(), t.Errorf, "placement is placed - it should not be")
	assert.That(
		pl.AtPoint() == game.Point{},
		t.Errorf, "placement reported at %#v - should be %#v", pl.AtPoint(), game.Point{})
}

func assertPlaced(t *testing.T, pl *game.HeadedPlacement, at game.Point) {
	assert.That(pl.Placed(), t.Errorf, "placement is not placed - it should be")
	assert.That(
		pl.AtPoint() == at,
		t.Errorf, "placement reported at %#v - should be %#v", pl.AtPoint(), at)
}

func assertNotHeaded(t *testing.T, pl *game.HeadedPlacement) {
	assert.That(!pl.Headed(), t.Errorf, "placement is headed - it should not be")
	assert.That(
		pl.Heading() == game.Point{},
		t.Errorf, "reported headint to %#v - should be %#v", pl.AtPoint(), game.Point{})
}

func assertHeaded(t *testing.T, pl *game.HeadedPlacement, to game.Point) {
	assert.That(pl.Headed(), t.Errorf, "placement is not headed - it should be")
	assert.That(
		pl.Heading() == to,
		t.Errorf, "heading reported at %#v - should be %#v", pl.Heading(), to)
}
