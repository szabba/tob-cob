// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game_test

import (
	"testing"

	"github.com/szabba/assert"
	"github.com/szabba/tob-cob/game"
)

func TestSpaceHasNonexistentPositionsByDefault(t *testing.T) {
	// given
	space := game.NewSpace()

	// when
	pos := space.At(game.P(13, 25))

	// then
	assert.That(!pos.Exists(), t.Errorf, "the position should not exist")
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

func TestDestroyedPositionShouldNotExist(t *testing.T) {
	// given
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

func TestItIsSafeToCreateAPositionTwice(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()

	// when
	p := catchPanic(func() {
		pos.Create()
	})

	// then
	assert.That(p == nil, t.Errorf, "a second call should not panic")
}

func TestItIsSafeToDestroyAPositionTwice(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))
	pos.Create()
	pos.Destroy()

	// when
	p := catchPanic(func() {
		pos.Destroy()
	})

	// then
	assert.That(p == nil, t.Errorf, "a second call should not panic")
}

func TestItIsSafeToCallDestroyBeforeCreate(t *testing.T) {
	// given
	space := game.NewSpace()
	pos := space.At(game.P(13, 25))

	// when
	p := catchPanic(func() {
		pos.Destroy()
	})

	// then
	assert.That(p == nil, t.Errorf, "the call should not panic")
}

func catchPanic(f func()) (p interface{}) {
	defer func() { p = recover() }()
	f()
	return p
}
