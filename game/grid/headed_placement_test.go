// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package grid_test

import (
	"testing"
	"time"

	"github.com/szabba/assert"

	"github.com/szabba/tob-cob/game/actions"
	"github.com/szabba/tob-cob/game/grid"
)

func TestZeroHeadedPlacementIsNotPlaced(t *testing.T) {
	// given
	// when
	zero := grid.HeadedPlacement{}

	// then
	assertNotPlaced(t, &zero)
}

func TestZeroHeadedPlacementIsNotHeaded(t *testing.T) {
	// given
	// when
	zero := grid.HeadedPlacement{}

	// then
	assertNotHeaded(t, &zero)
}

func TestZeroHeadedPlacementHasNoProgress(t *testing.T) {
	// given
	// when
	zero := grid.HeadedPlacement{}

	// then
	assert.That(zero.Progress() == 0, t.Errorf, "got progress %v - want %v", zero.Progress(), 0)
}

func TestHeadedPlacementCannotBePlacedAtPositionThatDoesNotExist(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))

	placement := grid.HeadedPlacement{}

	// when
	ok := placement.Place(pos)

	// then
	assert.That(!ok, t.Errorf, "placing should fail")
	assertNotPlaced(t, &placement)
	assert.That(!pos.Taken(), t.Errorf, "the position is taken - it should not be")
}

func TestHeadedPlacementCannotBePlacedAtPositionThatIsAlreadyTaken(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Take(grid.DummyTaker())

	placement := grid.HeadedPlacement{}

	// when
	ok := placement.Place(pos)

	// then
	assert.That(!ok, t.Errorf, "placing should fail")
	assertNotPlaced(t, &placement)
	assert.That(!pos.Taken(), t.Errorf, "the position is taken - it should not be")
}

func TestHeadedPlacementCanBePlaced(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()

	placement := grid.HeadedPlacement{}

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
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()

	placement := grid.HeadedPlacement{}

	// when
	action := placement.MoveTo(pos, time.Second)

	// then
	assert.That(
		action == actions.NoAction(),
		t.Errorf, "got action %#v, want %#v", action, actions.NoAction())
}

func TestPlacedHeadedPlacementCanMove(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()
	dst := space.At(grid.P(2, 4))
	dst.Create()

	placement := grid.HeadedPlacement{}
	placement.Place(pos)

	timeNeeded := time.Second
	action := placement.MoveTo(dst, timeNeeded)

	// when
	status := action.Run(timeNeeded)

	// then
	assert.That(
		status == actions.Done(0),
		t.Errorf, "got move status %#v - want %#v", status, actions.Done(0))
	assertPlaced(t, &placement, dst.AtPoint())
	assert.That(dst.Taken(), t.Errorf, "%#v should be taken, but is not", dst.AtPoint())
	assert.That(!pos.Taken(), t.Errorf, "%#v should not be taken, but it is", pos.AtPoint())
}

func TestPlacedHeadedPlacementHasCompleted(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()

	placement := grid.HeadedPlacement{}

	// when
	placement.Place(pos)

	// then
	assert.That(placement.Progress() == 1, t.Errorf, "got progress %f - want %f", placement.Progress(), 1)
}

func TestPlacedHeadedPlacementBecomesHeadedWhenAMoveStarts(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()
	dst := space.At(grid.P(2, 4))
	dst.Create()

	placement := grid.HeadedPlacement{}
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
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()
	dst := space.At(grid.P(2, 4))
	dst.Create()

	placement := grid.HeadedPlacement{}
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
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()

	placement := grid.HeadedPlacement{}

	// when
	placement.Place(pos)

	// then
	assertPlaced(t, &placement, pos.AtPoint())
	assertNotHeaded(t, &placement)
}

func TestHeadedPlacementLosesHeadingAfterBeingPlaced(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(2, 3))
	pos.Create()
	moveDst := space.At(grid.P(2, 4))
	moveDst.Create()
	dst := space.At(grid.P(1, 3))

	placement := grid.HeadedPlacement{}
	placement.Place(pos)

	timeNeeded := 4 * time.Second
	action := placement.MoveTo(moveDst, timeNeeded)
	action.Run(0)

	// when
	ok := placement.Place(dst)

	// then
	assert.That(ok, t.Errorf, "the placement should succeed - it failed")
	assertNotHeaded(t, &placement)
}

// TODO: Test (*HeadedPlacement).FollowPath

func TestFollowingAnEmptyPathCompletesImmediately(t *testing.T) {
	// given
	placement := grid.HeadedPlacement{}

	path := grid.Path{}
	stepTime := time.Second

	action := placement.FollowPath(path, stepTime)

	// when
	status := action.Run(0)

	// then
	want := actions.Done(0)
	assert.That(status == want, t.Errorf, "got status %#v, want %#v", status, want)
}

func TestFollowingASinglePositionPathCompletesImmediatelyWhenThePlacementIsAlreadyThere(t *testing.T) {
	// given
	space := grid.NewSpace()

	pos := space.At(grid.P(1, 2))
	pos.Create()

	placement := grid.HeadedPlacement{}
	placement.Place(pos)

	path, stepTime := grid.Path{pos}, time.Second

	action := placement.FollowPath(path, stepTime)

	// when
	status := action.Run(time.Second)

	// then
	want := actions.Done(time.Second)
	assert.That(status == want, t.Errorf, "got status %#v, want %#v", status, want)
}

func TestFollowingASinglePositionPathFailsImmediatelyWhenThePlacementIsAlreadyThere(t *testing.T) {
	// given
	space := grid.NewSpace()

	pos := space.At(grid.P(1, 2))
	pos.Create()

	placement := grid.HeadedPlacement{}

	path, stepTime := grid.Path{pos}, time.Second

	action := placement.FollowPath(path, stepTime)

	// when
	status := action.Run(time.Second)

	// then
	want := actions.Interrupted(time.Second)
	assert.That(status == want, t.Errorf, "got status %#v, want %#v", status, want)
}

func TestFollowingASingleStepPathCompletesAtTheStepTime(t *testing.T) {
	// given
	space := grid.NewSpace()

	src, dst := space.At(grid.P(1, 2)), space.At(grid.P(1, 3))
	src.Create()
	dst.Create()

	placement := grid.HeadedPlacement{}
	placement.Place(src)

	path, stepTime := grid.Path{src, dst}, time.Second

	action := placement.FollowPath(path, stepTime)

	// when
	status := action.Run(time.Second)

	// then
	want := actions.Done(0)
	assert.That(status == want, t.Errorf, "got status %#v, want %#v", status, want)
}

func TestFollowingASignleStepPathToCompletionPutsThePlacementAtRightPosition(t *testing.T) {
	// given
	space := grid.NewSpace()

	src, dst := space.At(grid.P(1, 2)), space.At(grid.P(1, 3))
	src.Create()
	dst.Create()

	placement := grid.HeadedPlacement{}
	placement.Place(src)

	path, stepTime := grid.Path{src, dst}, time.Second
	assertPathFromTo(assumption(t), path, src.AtPoint(), dst.AtPoint())

	action := placement.FollowPath(path, stepTime)

	// when
	action.Run(time.Second)

	// then
	want := dst.AtPoint()
	assert.That(placement.AtPoint() == want, t.Errorf, "placement at %#v, want %#v", placement.AtPoint(), want)
}

func TestFollowingATwoStepPathForHalfTheRequiredTimeLeavesTheActionPaused(t *testing.T) {
	// given
	space := grid.NewSpace()

	src, mid, dst := space.At(grid.P(1, 2)), space.At(grid.P(1, 3)), space.At(grid.P(1, 4))
	src.Create()
	mid.Create()
	dst.Create()

	placement := grid.HeadedPlacement{}
	placement.Place(src)

	path, stepTime := grid.Path{src, mid, dst}, time.Second
	assertPathFromTo(assumption(t), path, src.AtPoint(), dst.AtPoint())

	action := placement.FollowPath(path, stepTime)

	// when
	status := action.Run(time.Second)

	// then
	want := actions.Paused()
	assert.That(status == want, t.Errorf, "got action status %#v, want %#v", status, want)
}

func TestFollowingATwoStepPathForHalfTheRequiredTimeLeavesThePlacementAtTheMidpoint(t *testing.T) {
	// given
	space := grid.NewSpace()

	src, mid, dst := space.At(grid.P(1, 2)), space.At(grid.P(1, 3)), space.At(grid.P(1, 4))
	src.Create()
	mid.Create()
	dst.Create()

	placement := grid.HeadedPlacement{}
	placement.Place(src)

	path, stepTime := grid.Path{src, mid, dst}, time.Second
	assertPathFromTo(assumption(t), path, src.AtPoint(), dst.AtPoint())

	action := placement.FollowPath(path, stepTime)

	// when
	action.Run(time.Second)

	// then
	want := mid.AtPoint()
	assert.That(placement.AtPoint() == want, t.Errorf, "placement at %#v, want %#v", placement.AtPoint(), want)
}

func TestFollowingATwoStepPathForTheRequiredTimeLeavesActionDone(t *testing.T) {
	// given
	space := grid.NewSpace()

	src, mid, dst := space.At(grid.P(1, 2)), space.At(grid.P(1, 3)), space.At(grid.P(1, 4))
	src.Create()
	mid.Create()
	dst.Create()

	placement := grid.HeadedPlacement{}
	placement.Place(src)

	path, stepTime := grid.Path{src, mid, dst}, time.Second
	assertPathFromTo(assumption(t), path, src.AtPoint(), dst.AtPoint())

	action := placement.FollowPath(path, stepTime)

	// when
	status := action.Run(2 * time.Second)

	// then
	// then
	want := actions.Done(0)
	assert.That(status == want, t.Errorf, "got action status %#v, want %#v", status, want)
}

func TestFollowingATwoStepPathForTheRequiredTimeLeavesThePlacementAtTheDestination(t *testing.T) {
	// given
	space := grid.NewSpace()

	src, mid, dst := space.At(grid.P(1, 2)), space.At(grid.P(1, 3)), space.At(grid.P(1, 4))
	src.Create()
	mid.Create()
	dst.Create()

	placement := grid.HeadedPlacement{}
	placement.Place(src)

	path, stepTime := grid.Path{src, mid, dst}, time.Second
	assertPathFromTo(assumption(t), path, src.AtPoint(), dst.AtPoint())

	action := placement.FollowPath(path, stepTime)

	// when
	action.Run(2 * time.Second)

	// then
	want := dst.AtPoint()
	assert.That(placement.AtPoint() == want, t.Errorf, "placement at %#v, want %#v", placement.AtPoint(), want)
}

// TODO: Better name?
func assumption(t *testing.T) assert.ErrorFunc {
	return func(msg string, args ...interface{}) {
		t.Fatalf("assumption violated: "+msg, args...)
	}
}

func assertNotPlaced(t *testing.T, pl *grid.HeadedPlacement) {
	assert.That(!pl.Placed(), t.Errorf, "placement is placed - it should not be")
	assert.That(
		pl.AtPoint() == grid.Point{},
		t.Errorf, "placement reported at %#v - should be %#v", pl.AtPoint(), grid.Point{})
}

func assertPlaced(t *testing.T, pl *grid.HeadedPlacement, at grid.Point) {
	assert.That(pl.Placed(), t.Errorf, "placement is not placed - it should be")
	assert.That(
		pl.AtPoint() == at,
		t.Errorf, "placement reported at %#v - should be %#v", pl.AtPoint(), at)
}

func assertNotHeaded(t *testing.T, pl *grid.HeadedPlacement) {
	assert.That(!pl.Headed(), t.Errorf, "placement is headed - it should not be")
	assert.That(
		pl.Heading() == grid.Point{},
		t.Errorf, "reported heading to %#v - should be %#v", pl.AtPoint(), grid.Point{})
}

func assertHeaded(t *testing.T, pl *grid.HeadedPlacement, to grid.Point) {
	assert.That(pl.Headed(), t.Errorf, "placement is not headed - it should be")
	assert.That(
		pl.Heading() == to,
		t.Errorf, "heading reported at %#v - should be %#v", pl.Heading(), to)
}
