// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package grid_test

import (
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/tob-cob/game/grid"
)

func TestEmptyPathIsNotViable(t *testing.T) {
	// given
	space := grid.NewSpace()
	finder := grid.NewPathFinder(space)

	path := grid.Path{}

	// when
	viable := finder.IsViable(path)

	// then
	assert.That(!viable, t.Errorf, "path %#v should not be viable", path)
}

func TestPathOfOneExistingPointIsViable(t *testing.T) {
	// given
	space := grid.NewSpace()
	space.At(grid.P(3, 4)).Create()
	finder := grid.NewPathFinder(space)

	path := pathOf(space, grid.P(3, 4))

	// when
	viable := finder.IsViable(path)

	// then
	assert.That(viable, t.Errorf, "path %#v should be viable", path)
}

func TestPathOfOnePointThatDoesNotExistIsUnviable(t *testing.T) {
	// given
	space := grid.NewSpace()
	finder := grid.NewPathFinder(space)

	path := pathOf(space, grid.P(3, 4))

	// when
	viable := finder.IsViable(path)

	// then
	assert.That(!viable, t.Errorf, "path %#v should not be viable", path)
}

func TestPathWithJumpsIsNotViable(t *testing.T) {
	// given
	space := grid.NewSpace()
	finder := grid.NewPathFinder(space)

	points := []grid.Point{grid.P(3, 4), grid.P(3, 6)}
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
	space := grid.NewSpace()
	finder := grid.NewPathFinder(space)

	points := []grid.Point{grid.P(3, 4), grid.P(4, 5)}
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
	space := grid.NewSpace()
	space.At(grid.P(3, 4)).Create()
	finder := grid.NewPathFinder(space)

	src := space.At(grid.P(3, 4))
	dst := space.At(grid.P(3, 4))

	// when
	path, ok := finder.FindPath(src, dst)

	// then
	assert.That(ok, t.Errorf, "path search failed")
	assert.That(finder.IsViable(path), t.Errorf, "found unviable path %#v", path)
	assertPathFromTo(t.Errorf, path, grid.P(3, 4), grid.P(3, 4))
}

func TestPathIsNotFoundBetweenDisconnectedPoints(t *testing.T) {
	// given
	space := grid.NewSpace()
	space.At(grid.P(3, 4)).Create()
	space.At(grid.P(3, 6)).Create()
	finder := grid.NewPathFinder(space)

	src := space.At(grid.P(3, 4))
	dst := space.At(grid.P(3, 6))

	// when
	path, ok := finder.FindPath(src, dst)

	// then
	assert.That(!ok, t.Errorf, "path search did not fail")
	assert.That(finder.IsViable(path), t.Errorf, "found unviable path %#v", path)
	assertPathFromTo(t.Errorf, path, grid.P(3, 4), grid.P(3, 4))
}

func TestPathIsFoundBetweenConnectedPoints(t *testing.T) {
	// given
	space := grid.NewSpace()
	space.At(grid.P(3, 4)).Create()
	space.At(grid.P(3, 5)).Create()
	space.At(grid.P(4, 5)).Create()
	finder := grid.NewPathFinder(space)

	src := space.At(grid.P(3, 4))
	dst := space.At(grid.P(4, 5))

	// when
	path, ok := finder.FindPath(src, dst)

	// then
	assert.That(ok, t.Errorf, "path search failed")
	assert.That(finder.IsViable(path), t.Errorf, "found unviable path %#v", path)
	assertPathFromTo(t.Errorf, path, grid.P(3, 4), grid.P(4, 5))
}

func TestFoundPathAvoidsTakenPositions(t *testing.T) {
	// given
	space := grid.NewSpace()
	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			space.At(grid.P(y, x)).Create()
		}
	}
	taker := grid.OnePosTaker{}
	taken := space.At(grid.P(0, 1))
	taken.Take(&taker)

	finder := grid.NewPathFinder(space)

	src := space.At(grid.P(0, 0))
	dst := space.At(grid.P(0, 2))

	// when
	path, ok := finder.FindPath(src, dst)

	// then
	assert.That(ok, t.Errorf, "path search failed")
	assert.That(finder.IsViable(path), t.Errorf, "found unviable path %#v", path)
	assertPathFromTo(t.Errorf, path, src.AtPoint(), dst.AtPoint())
	for i, pos := range path {
		assert.That(pos != taken, t.Errorf, "at index %d: path includes taken position %#v", i, taken)
	}
}

func assertPathFromTo(onErr assert.ErrorFunc, path grid.Path, from, to grid.Point) {
	if len(path) == 0 {
		onErr("path is empty")
		return
	}

	last := len(path) - 1
	assert.That(path[0].AtPoint() == from, onErr, "got path start %#v - want %#v", path[0], from)
	assert.That(path[last].AtPoint() == to, onErr, "got path end %#v - want %#v", path[last], to)

}

func pathOf(space *grid.Space, pts ...grid.Point) []grid.Position {
	path := make([]grid.Position, len(pts))
	for i, pt := range pts {
		path[i] = space.At(pt)
	}
	return path
}
