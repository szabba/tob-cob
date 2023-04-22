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

func TestPositionsFromTheSameSpaceAtTheSamePointAreEqual(t *testing.T) {
	// given
	space := grid.NewSpace()
	first := space.At(grid.P(13, 25))

	// when
	second := space.At(grid.P(13, 25))

	// then
	assert.That(
		first == second,
		t.Errorf, "%#v and %#v are not the same", first, second)
}

func TestPositionsFromTheTwoSpacesAtTheSamePointAreNotEqual(t *testing.T) {
	// given
	space := grid.NewSpace()
	otherSpace := grid.NewSpace()
	first := space.At(grid.P(13, 25))

	// when
	second := otherSpace.At(grid.P(13, 25))

	// then
	assert.That(
		first != second,
		t.Errorf, "different spaces return the same position %#v", first)
}
func TestPositionIsAtItsPoint(t *testing.T) {
	// given
	space := grid.NewSpace()

	// when
	pos := space.At(grid.P(13, 25))

	// then
	assert.That(
		pos.AtPoint() == grid.P(13, 25),
		t.Errorf, "position at %#v, want %#v", pos.AtPoint(), grid.P(13, 25))
}

func TestPositionDoesNotExistByDefault(t *testing.T) {
	// given
	space := grid.NewSpace()

	// when
	pos := space.At(grid.P(13, 25))

	// then
	assert.That(!pos.Exists(), t.Errorf, "the position should not exist")
}

func TestPositionThatDoesNotExist(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))

	// when
	ok := pos.Create()

	// then
	assert.That(ok, t.Errorf, "creating the position should succeed")
}

func TestPoisitionExistsOnceItHasBeenCreated(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))

	// when
	pos.Create()

	// then
	assert.That(pos.Exists(), t.Errorf, "the position should exist")
}

func TestOnlyThePositionCreatedExists(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()

	// when
	otherPos := space.At(grid.P(25, -13))

	// then
	assert.That(!otherPos.Exists(), t.Errorf, "the other position should not exist")
}

func TestACreatedPositionCanBeDestroyed(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()

	// when
	ok := pos.Destroy()

	// then
	assert.That(ok, t.Errorf, "destroying the position should succeed")
}

func TestDestroyedPositionShouldNotExist(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()

	// when
	pos.Destroy()

	// then
	assert.That(!pos.Exists(), t.Errorf, "the position should not exist")
}

func TestDestroyedPositionCanBeRecreated(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Destroy()

	// when
	pos.Create()

	// then
	assert.That(pos.Exists(), t.Errorf, "the position should exist")
}

func TestAPositionCannotBeCreatedTwiceInARow(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()

	// when
	ok := pos.Create()

	// then
	assert.That(!ok, t.Errorf, "creating a position twice should fail")
}

func TestAPositionCannotBeDestroyedTwiceInARow(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Destroy()

	// when
	ok := pos.Destroy()

	// then
	assert.That(!ok, t.Errorf, "destroying the position should fail")
}

func TestAPositionThatWasNotCreatedCannotBeDestroyed(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))

	// when
	ok := pos.Destroy()

	// then
	assert.That(!ok, t.Errorf, "destroying the position should fail")
}

func TestAPositionThatDoesNotExistIsNotTaken(t *testing.T) {
	// given
	space := grid.NewSpace()

	// when
	pos := space.At(grid.P(13, 25))

	// then
	assert.That(!pos.Taken(), t.Errorf, "the position should not be taken")
}

func TestAPositionThatWasCreatedIsNotTakenByDefault(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))

	// when
	pos.Create()

	// then
	assert.That(!pos.Taken(), t.Errorf, "the position should not be taken")
}

func TestAPositionThatDoesNotExistCannotBeTaken(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))

	// when
	ok := pos.Take(grid.DummyTaker())

	// then
	assert.That(!ok, t.Errorf, "taking the pos should fail")
}

func TestAPositionThatExistCanBeTaken(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()

	// when
	ok := pos.Take(grid.DummyTaker())

	// then
	assert.That(ok, t.Errorf, "taking the position should succeed")
}

func TestAPositionCannotBeTakenTwice(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	ok := pos.Take(grid.DummyTaker())

	// then
	assert.That(!ok, t.Errorf, "taking the position for a second time should fail")
}

func TestAPositionThatWasTakenSaysSo(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	taken := pos.Taken()

	// then
	assert.That(taken, t.Errorf, "the position was taken but says it was not")
}

func TestAPositionCreatedAgainRemainsTaken(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	pos.Create()

	// then
	taken := pos.Taken()
	assert.That(taken, t.Errorf, "the position was taken but says it was not")
}

func TestATakenPositionCannotBeCreated(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	ok := pos.Create()

	// then
	assert.That(!ok, t.Errorf, "creating the position should fail")
}

func TestAPositionRemainsTakenWhenAnAttemptIsMadeToCreateIt(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	pos.Create()

	// then
	assert.That(pos.Taken(), t.Errorf, "the position should still be taken")
}

func TestATakenPositionCannotBeDestroyed(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	ok := pos.Destroy()

	// then
	assert.That(!ok, t.Errorf, "destroying the position should fail")
}

func TestAPositionThatDoesNotExistCannotBeFreed(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))

	// when
	ok := pos.Free()

	// then
	assert.That(!ok, t.Errorf, "freeing the position should fail")
}

func TestAPositionThatIsNotTakenCannotBeFreed(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()

	// when
	ok := pos.Free()

	// then
	assert.That(!ok, t.Errorf, "freeing the position should fail")
}

func TestAPositionThatIsTakenCanBeFreed(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	ok := pos.Free()

	// then
	assert.That(ok, t.Errorf, "freeing the position should succeed")
}

func TestAFreedPositionIsNotTaken(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	pos.Free()

	// then
	assert.That(!pos.Taken(), t.Errorf, "position should not be taken")
}

func TestSpaceTakerIsNotifiedWhenItIsLetOntoAPosition(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()

	taker := &RecordingSpaceTaker{}

	// when
	pos.Take(taker)

	// then
	expectedCall := RecordedSpaceTakerCall{
		Method:   "LetOnto",
		Position: pos,
	}
	assert.That(
		len(taker.Calls) == 1,
		t.Fatalf, "got %d calls, want %d", len(taker.Calls), 1)
	assert.That(
		taker.Calls[0] == expectedCall,
		t.Fatalf, "got %d call %#v, wan %#v", 0, taker.Calls[0], expectedCall)
}

func TestSpaceTakerIsNotifiedWhenItIsForcedOutOfAPosition(t *testing.T) {
	// given
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))
	pos.Create()

	taker := &RecordingSpaceTaker{}
	pos.Take(taker)

	// when
	pos.Free()

	// then
	expectedCall := RecordedSpaceTakerCall{
		Method:   "ForceOff",
		Position: pos,
	}
	assert.That(
		len(taker.Calls) == 2,
		t.Fatalf, "got %d calls, want %d", len(taker.Calls), 2)
	assert.That(
		taker.Calls[1] == expectedCall,
		t.Fatalf, "got %d call %#v, wan %#v", 1, taker.Calls[1], expectedCall)
}

func TestActionTakingAFreePositionSucceeds(t *testing.T) {
	// given
	action, pos, _ := setUpActionTakingPosition()
	pos.Create()

	// when
	status := action.Run(time.Second)

	// then
	assert.That(
		status == actions.Done(time.Second),
		t.Errorf, "got status %#v, want %#v", status, actions.Done(time.Second))
}

func TestActionTakingAFreePositionCallsTheSpaceTaker(t *testing.T) {
	// given
	action, pos, taker := setUpActionTakingPosition()
	pos.Create()

	// when
	action.Run(time.Second)

	// then
	expectedCall := RecordedSpaceTakerCall{Method: "LetOnto", Position: pos}
	assert.That(
		len(taker.Calls) == 1,
		t.Fatalf, "got %d space taker calls, want %d", len(taker.Calls), 1)
	assert.That(
		taker.Calls[0] == expectedCall,
		t.Errorf, "got %d call %#v, want %#v", taker.Calls[0], expectedCall)
}

func TestActionTakingAFreePositionTakesIt(t *testing.T) {
	// given
	action, pos, _ := setUpActionTakingPosition()
	pos.Create()

	// when
	action.Run(time.Second)

	// then
	assert.That(
		pos.Taken(),
		t.Errorf, "position at %#v should be taken", pos.AtPoint())
}

func TestActionTakingATakenPositionFails(t *testing.T) {
	// given
	action, pos, _ := setUpActionTakingPosition()
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	status := action.Run(time.Second)

	// then
	assert.That(
		status == actions.Interrupted(time.Second),
		t.Errorf, "got status %#v, want %#v", status, actions.Interrupted(time.Second))
}

func TestActionTakingATakenPositionLeavesThePositionStateAsIs(t *testing.T) {
	// given
	action, pos, _ := setUpActionTakingPosition()
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	action.Run(time.Second)

	// then
	assert.That(
		pos.Taken(),
		t.Errorf, "the position %#v should be taken", pos)
}

func TestActionTakingATakenPositionDoesNotCallTheSpaceTaker(t *testing.T) {
	// given
	action, pos, taker := setUpActionTakingPosition()
	pos.Create()
	pos.Take(grid.DummyTaker())

	// when
	action.Run(time.Second)

	// then
	assert.That(
		len(taker.Calls) == 0,
		t.Fatalf, "got %d space taker calls, want %d", len(taker.Calls), 0)
}

func TestActionTakingANonExistentPositionFails(t *testing.T) {
	// given
	action, _, _ := setUpActionTakingPosition()

	// when
	status := action.Run(time.Second)

	// then
	assert.That(
		status == actions.Interrupted(time.Second),
		t.Errorf, "got status %#v, want %#v", status, actions.Interrupted(time.Second))
}

func TestActionTakingANonExistentPositionLeavesThePositionStateAsIs(t *testing.T) {
	// given
	action, pos, _ := setUpActionTakingPosition()

	// when
	action.Run(time.Second)

	// then
	assert.That(
		!pos.Exists(),
		t.Errorf, "the position %#v should not exist", pos)
}

func TestActionTakingANonExistentPositionDoesNotCallTheSpaceTaker(t *testing.T) {
	// given
	action, _, taker := setUpActionTakingPosition()

	// when
	action.Run(time.Second)

	// then
	assert.That(
		len(taker.Calls) == 0,
		t.Fatalf, "got %d space taker calls, want %d", len(taker.Calls), 0)
}

func TestEmptySpaceHasZeroMin(t *testing.T) {
	// given
	// when
	space := grid.NewSpace()

	// then
	assert.That(space.Min() == grid.Point{}, t.Errorf, "got %#v - want %#v", space.Min(), grid.Point{})
}

func TestEmptySpaceHasZeroMax(t *testing.T) {
	// given
	// when
	space := grid.NewSpace()

	// then
	assert.That(space.Max() == grid.Point{}, t.Errorf, "got %#v - want %#v", space.Max(), grid.Point{})
}

func TestSpaceWithOriginHasZeroMin(t *testing.T) {
	// given
	space := grid.NewSpace()

	// when
	space.At(grid.P(0, 0)).Create()

	// then
	assert.That(space.Min() == grid.Point{}, t.Errorf, "got %#v - want %#v", space.Min(), grid.Point{})
}

func TestSpaceWithOriginHasZeroMax(t *testing.T) {
	// given
	space := grid.NewSpace()

	// when
	space.At(grid.P(0, 0)).Create()

	// then
	assert.That(space.Max() == grid.Point{}, t.Errorf, "got %#v - want %#v", space.Max(), grid.Point{})
}

func TestSpaceWithOnePointHasItAsMin(t *testing.T) {
	// given
	space := grid.NewSpace()
	pt := grid.P(3, 7)

	// when
	space.At(pt).Create()

	// then
	assert.That(space.Min() == pt, t.Errorf, "got %#v - want %#v", space.Min(), pt)
}

func TestSpaceWithOnePointHasIsAsMax(t *testing.T) {
	// given
	space := grid.NewSpace()
	pt := grid.P(3, 7)

	// when
	space.At(pt).Create()

	// then
	assert.That(space.Max() == pt, t.Errorf, "got %#v - want %#v", space.Max(), pt)
}

func TestSpaceWithThreePointsInARowHasTheLeftmostAsMin(t *testing.T) {
	// given
	space := grid.NewSpace()
	left, mid, right := grid.P(3, 7), grid.P(3, 8), grid.P(3, 10)

	// when
	space.At(left).Create()
	space.At(mid).Create()
	space.At(right).Create()

	// then
	assert.That(space.Min() == left, t.Errorf, "got %#v - want %#v", space.Min(), left)
}

func TestSpaceWithThreePointsInARowHasTheRightmostAsMax(t *testing.T) {
	// given
	space := grid.NewSpace()
	left, mid, right := grid.P(3, 7), grid.P(3, 8), grid.P(3, 10)

	// when
	space.At(left).Create()
	space.At(mid).Create()
	space.At(right).Create()

	// then
	assert.That(space.Max() == right, t.Errorf, "got %#v - want %#v", space.Max(), right)
}

func TestSpaceWithThreePointsInARowHasLeftmostAsMinWhenPositionsAreCreatedOutOfOrder(t *testing.T) {
	// given
	space := grid.NewSpace()
	left, mid, right := grid.P(3, 7), grid.P(3, 8), grid.P(3, 10)

	// when
	space.At(mid).Create()
	space.At(right).Create()
	space.At(left).Create()

	// then
	assert.That(space.Max() == right, t.Errorf, "got %#v - want %#v", space.Max(), right)
}

func TestSpaceWithThreePointsInARowHasRightmostAsMaxWhenPositionsAreCreatedOutOfOrder(t *testing.T) {
	// given
	space := grid.NewSpace()
	left, mid, right := grid.P(3, 7), grid.P(3, 8), grid.P(3, 10)

	// when
	space.At(left).Create()
	space.At(right).Create()
	space.At(mid).Create()

	// then
	assert.That(space.Max() == right, t.Errorf, "got %#v - want %#v", space.Max(), right)
}

func TestSpaceWithThreePointsInARowHasMiddleOneAsMinAfterTheLeftmostOneIsRemoved(t *testing.T) {
	// given
	space := grid.NewSpace()
	left, mid, right := grid.P(3, 7), grid.P(3, 8), grid.P(3, 10)

	space.At(left).Create()
	space.At(mid).Create()
	space.At(right).Create()

	// when
	space.At(left).Destroy()

	// then
	assert.That(space.Min() == mid, t.Errorf, "got %#v - want %#v", space.Min(), mid)
}

func TestSpaceWithThreePointsInARowHasMiddleOneAsMaxAfterTheRightmostOneIsRemoved(t *testing.T) {
	// given
	space := grid.NewSpace()
	left, mid, right := grid.P(3, 7), grid.P(3, 8), grid.P(3, 10)

	space.At(left).Create()
	space.At(mid).Create()
	space.At(right).Create()

	// when
	space.At(right).Destroy()

	// then
	assert.That(space.Max() == mid, t.Errorf, "got %#v - want %#v", space.Max(), mid)
}

func TestSpaceWithThreePointsInAColumnHasTheBottomOneAsMin(t *testing.T) {
	// given
	space := grid.NewSpace()
	bottom, mid, top := grid.P(3, 7), grid.P(4, 7), grid.P(6, 7)

	// when
	space.At(bottom).Create()
	space.At(mid).Create()
	space.At(top).Create()

	// then
	assert.That(space.Min() == bottom, t.Errorf, "got %#v - want %#v", space.Min(), bottom)
}

func TestSpaceWithThreePointsInAColumnHasTheTopOneAsMax(t *testing.T) {
	// given
	space := grid.NewSpace()
	bottom, mid, top := grid.P(3, 7), grid.P(4, 7), grid.P(6, 7)

	// when
	space.At(bottom).Create()
	space.At(mid).Create()
	space.At(top).Create()

	// then
	assert.That(space.Max() == top, t.Errorf, "got %#v - want %#v", space.Max(), top)
}

func TestSpaceWithThreePointsInAColumnHasBottomOneAsMinWhenPositionsAreCreatedOutOfOrder(t *testing.T) {
	// given
	space := grid.NewSpace()
	bottom, mid, top := grid.P(3, 7), grid.P(4, 7), grid.P(6, 7)

	// when
	space.At(mid).Create()
	space.At(top).Create()
	space.At(bottom).Create()

	// then
	assert.That(space.Min() == bottom, t.Errorf, "got %#v - want %#v", space.Min(), bottom)
}

func TestSpaceWithThreePointsInAColumnHasTopOneAsMaxWhenPositionsAreCreatedOutOfOrder(t *testing.T) {
	// given
	space := grid.NewSpace()
	bottom, mid, top := grid.P(3, 7), grid.P(4, 7), grid.P(6, 7)

	// when
	space.At(top).Create()
	space.At(bottom).Create()
	space.At(mid).Create()

	// then
	assert.That(space.Max() == top, t.Errorf, "got %#v - want %#v", space.Max(), top)
}

func TestSpaceWithThreePointsInAColumnHasMiddleOneAsMinAfterTheBottomOneIsRemoved(t *testing.T) {
	// given
	space := grid.NewSpace()
	bottom, mid, top := grid.P(3, 7), grid.P(4, 7), grid.P(6, 7)

	space.At(bottom).Create()
	space.At(mid).Create()
	space.At(top).Create()

	// when
	space.At(bottom).Destroy()

	// then
	assert.That(space.Min() == mid, t.Errorf, "got %#v - want %#v", space.Min(), mid)
}

func TestSpaceWithThreePointsInAColumnHasMiddleOneAsMaxAfterTheTopOneIsRemoved(t *testing.T) {
	// given
	space := grid.NewSpace()
	bottom, mid, top := grid.P(3, 7), grid.P(4, 7), grid.P(6, 7)

	space.At(bottom).Create()
	space.At(mid).Create()
	space.At(top).Create()

	// when
	space.At(top).Destroy()

	// then
	assert.That(space.Max() == mid, t.Errorf, "got %#v - want %#v", space.Max(), mid)
}

func setUpActionTakingPosition() (actions.Action, grid.Position, *RecordingSpaceTaker) {
	space := grid.NewSpace()
	pos := space.At(grid.P(13, 25))

	var taker RecordingSpaceTaker

	action := grid.TakePosition(pos, &taker)

	return action, pos, &taker
}

type RecordingSpaceTaker struct {
	Calls []RecordedSpaceTakerCall
}

var _ grid.SpaceTaker = &RecordingSpaceTaker{}

func (taker *RecordingSpaceTaker) LetOnto(pos grid.Position) {
	call := RecordedSpaceTakerCall{
		Method:   "LetOnto",
		Position: pos,
	}
	taker.Calls = append(taker.Calls, call)
}

func (taker *RecordingSpaceTaker) ForceOff(pos grid.Position) {
	call := RecordedSpaceTakerCall{
		Method:   "ForceOff",
		Position: pos,
	}
	taker.Calls = append(taker.Calls, call)
}

type RecordedSpaceTakerCall struct {
	Method   string
	Position grid.Position
}
