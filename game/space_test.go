// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game_test

import (
	"testing"

	"github.com/szabba/assert"
	"github.com/szabba/tob-cob/game"
)

func TestPositionIsAtItsPoint(t *testing.T) {
	// given
	space := game.NewSpace()

	// when
	pos := space.At(game.P(13, 25))

	// then
	assert.That(
		pos.AtPoint() == game.P(13, 25),
		t.Errorf, "position at %#v, want %#v", pos.AtPoint(), game.P(13, 25))
}

func TestPositionDoesNotExistByDefault(t *testing.T) {
	// given
	space := game.NewSpace()

	// when
	pos := space.At(game.P(13, 25))

	// then
	assert.That(!pos.Exists(), t.Errorf, "the position should not exist")
}

func TestPositionThatDoesNotExist(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))

	// when
	ok := pos.Create()

	// then
	assert.That(ok, t.Errorf, "creating the position should succeed")
}

func TestPoisitionExistsOnceItHasBeenCreated(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))

	// when
	pos.Create()

	// then
	assert.That(pos.Exists(), t.Errorf, "the position should exist")
}

func TestOnlyThePositionCreatedExists(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()

	// when
	otherPos := space.At(game.P(25, -13))

	// then
	assert.That(!otherPos.Exists(), t.Errorf, "the other position should not exist")
}

func TestACreatedPositionCanBeDestroyed(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()

	// when
	ok := pos.Destroy()

	// then
	assert.That(ok, t.Errorf, "destroying the position should succeed")
}

func TestDestroyedPositionShouldNotExist(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()

	// when
	pos.Destroy()

	// then
	assert.That(!pos.Exists(), t.Errorf, "the position should not exist")
}

func TestDestroyedPositionCanBeRecreated(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()
	pos.Destroy()

	// when
	pos.Create()

	// then
	assert.That(pos.Exists(), t.Errorf, "the position should exist")
}

func TestAPositionCannotBeCreatedTwiceInARow(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()

	// when
	ok := pos.Create()

	// then
	assert.That(!ok, t.Errorf, "creating a position twice should fail")
}

func TestAPositionCannotBeDestroyedTwiceInARow(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Destroy()

	// when
	ok := pos.Destroy()

	// then
	assert.That(!ok, t.Errorf, "destroying the position should fail")
}

func TestAPositionThatWasNotCreatedCannotBeDestroyed(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))

	// when
	ok := pos.Destroy()

	// then
	assert.That(!ok, t.Errorf, "destroying the position should fail")
}

func TestAPositionThatDoesNotExistIsNotTaken(t *testing.T) {
	// given
	space := game.NewSpace()

	// when
	pos := space.At(game.P(13, 25))

	// then
	assert.That(!pos.Taken(), t.Errorf, "the position should not be taken")
}

func TestAPositionThatWasCreatedIsNotTakenByDefault(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))

	// when
	pos.Create()

	// then
	assert.That(!pos.Taken(), t.Errorf, "the position should not be taken")
}

func TestAPositionThatDoesNotExistCannotBeTaken(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))

	// when
	ok := pos.Take()

	// then
	assert.That(!ok, t.Errorf, "taking the pos should fail")
}

func TestAPositionThatExistCanBeTaken(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()

	// when
	ok := pos.Take()

	// then
	assert.That(ok, t.Errorf, "taking the position should succeed")
}

func TestAPositionCannotBeTakenTwice(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()
	pos.Take()

	// when
	ok := pos.Take()

	// then
	assert.That(!ok, t.Errorf, "taking the position for a second time should fail")
}

func TestAPositionThatWasTakenSaysSo(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()
	pos.Take()

	// when
	taken := pos.Taken()

	// then
	assert.That(taken, t.Errorf, "the position was taken but says it was not")
}

func TestAPositionCreatedAgainRemainsTaken(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()
	pos.Take()

	// when
	pos.Create()

	// then
	taken := pos.Taken()
	assert.That(taken, t.Errorf, "the position was taken but says it was not")
}

func TestATakenPositionCannotBeDestroyed(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()
	pos.Take()

	// when
	ok := pos.Destroy()

	// then
	assert.That(!ok, t.Errorf, "desroying the position should fail")
}

func TestAPositionThatDoesNotExistCannotBeFreed(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))

	// when
	ok := pos.Free()

	// then
	assert.That(!ok, t.Errorf, "freeing the position should fail")
}

func TestAPositionThatIsNotTakenCannotBeFreed(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()

	// when
	ok := pos.Free()

	// then
	assert.That(!ok, t.Errorf, "freeing the position should fail")
}

func TestAPositionThatIsTakenCanBeFreed(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()
	pos.Take()

	// when
	ok := pos.Free()

	// then
	assert.That(ok, t.Errorf, "freeing the position should succeed")
}

func TestAFreedPositionIsNotTaken(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()
	pos.Take()

	// when
	pos.Free()

	// then
	assert.That(!pos.Taken(), t.Errorf, "position should not be taken")
}
