// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package game_test

import (
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/tob-cob/game"
)

func TestEmptyPathIsNotViable(t *testing.T) {
	// given
	space := game.NewSpace()
	finder := game.NewPathFinder(space)

	path := game.Path{}

	// when
	viable := finder.IsViable(path)

	// then
	assert.That(!viable, t.Errorf, "path %#v should not be viable", path)
}

func TestPathOfOneExistingPointIsViable(t *testing.T) {
	// given
	space := game.NewSpace()
	space.At(game.P(3, 4)).Create()
	finder := game.NewPathFinder(space)

	path := pathOf(space, game.P(3, 4))

	// when
	viable := finder.IsViable(path)

	// then
	assert.That(viable, t.Errorf, "path %#v should be viable", path)
}

func TestPathOfOnePointThatDoesNotExistIsUnviable(t *testing.T) {
	// given
	space := game.NewSpace()
	finder := game.NewPathFinder(space)

	path := pathOf(space, game.P(3, 4))

	// when
	viable := finder.IsViable(path)

	// then
	assert.That(!viable, t.Errorf, "path %#v should not be viable", path)
}

func TestPathWithJumpsIsNotViable(t *testing.T) {
	// given
	space := game.NewSpace()
	finder := game.NewPathFinder(space)

	points := []game.Point{game.P(3, 4), game.P(3, 6)}
	for _, pt := range points {
		space.At(pt).Create()
	}

	path := pathOf(space, points...)

	// when
	viable := finder.IsViable(path)

	// then
	assert.That(!viable, t.Errorf, "path %#v should not be viable", path)
}

func TestPathWithDiagonalsIsNotViable(t *testing.T) {
	// given
	space := game.NewSpace()
	finder := game.NewPathFinder(space)

	points := []game.Point{game.P(3, 4), game.P(4, 5)}
	for _, pt := range points {
		space.At(pt).Create()
	}

	path := pathOf(space, points...)

	// when
	viable := finder.IsViable(path)

	// then
	assert.That(!viable, t.Errorf, "path %#v should not be viable", path)
}

func TestPathFromPointToItselfContainsOnlyIt(t *testing.T) {
	// given
	space := game.NewSpace()
	space.At(game.P(3, 4)).Create()
	finder := game.NewPathFinder(space)

	src := space.At(game.P(3, 4))
	dst := space.At(game.P(3, 4))

	// when
	path, ok := finder.FindPath(src, dst)

	// then
	assert.That(ok, t.Errorf, "path search failed")
	assert.That(finder.IsViable(path), t.Errorf, "found unviable path %#v", path)
	assertPathFromTo(t.Errorf, path, game.P(3, 4), game.P(3, 4))
}

func TestPathIsNotFoundBetweenDisconnectedPoints(t *testing.T) {
	// given
	space := game.NewSpace()
	space.At(game.P(3, 4)).Create()
	space.At(game.P(3, 6)).Create()
	finder := game.NewPathFinder(space)

	src := space.At(game.P(3, 4))
	dst := space.At(game.P(3, 6))

	// when
	path, ok := finder.FindPath(src, dst)

	// then
	assert.That(!ok, t.Errorf, "path search did not fail")
	assert.That(finder.IsViable(path), t.Errorf, "found unviable path %#v", path)
	assertPathFromTo(t.Errorf, path, game.P(3, 4), game.P(3, 4))
}

func TestPathIsFoundBetweenConnectedPoints(t *testing.T) {
	// given
	space := game.NewSpace()
	space.At(game.P(3, 4)).Create()
	space.At(game.P(3, 5)).Create()
	space.At(game.P(4, 5)).Create()
	finder := game.NewPathFinder(space)

	src := space.At(game.P(3, 4))
	dst := space.At(game.P(4, 5))

	// when
	path, ok := finder.FindPath(src, dst)

	// then
	assert.That(ok, t.Errorf, "path search failed")
	assert.That(finder.IsViable(path), t.Errorf, "found unviable path %#v", path)
	assertPathFromTo(t.Errorf, path, game.P(3, 4), game.P(4, 5))
}

func assertPathFromTo(onErr assert.ErrorFunc, path game.Path, from, to game.Point) {
	if len(path) == 0 {
		onErr("path is empty")
		return
	}

	last := len(path) - 1
	assert.That(path[0].AtPoint() == from, onErr, "got path start %#v - want %#v", path[0], from)
	assert.That(path[last].AtPoint() == to, onErr, "got path end %#v - want %#v", path[last], to)

}

func pathOf(space *game.Space, pts ...game.Point) []game.Position {
	path := make([]game.Position, len(pts))
	for i, pt := range pts {
		path[i] = space.At(pt)
	}
	return path
}
