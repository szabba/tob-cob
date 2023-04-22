// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package grid

import (
	"math"

	"github.com/fzipp/astar"
)

// A Path is a sequence of positions.
type Path []Position

// A PathFinder finds paths from one point to another.
type PathFinder struct {
	space *Space
}

// NewPathFinder creates a path finder that searches for path through the specified space.
//
// The path finder will be sensitive to what positions do and do not exist in the space.
// The space can be modified after the path finder is created - it will be aware of the updates.
func NewPathFinder(space *Space) PathFinder {
	return PathFinder{space}
}

// IsViable validates a path.
//
// It will be false if the path contains positions from a different space.
// It will be false if the path contains positions that cannot be occupied.
// It will be false if the path contains two consecutive positions that are not neighbours.
func (pf PathFinder) IsViable(path Path) bool {
	viable := true
	viable = viable && len(path) > 0
	for _, pos := range path {
		inSpace := pf.space.At(pos.AtPoint()) == pos
		viable = viable && inSpace && pos.Exists()
	}
	for i := 1; i < len(path); i++ {
		prev, next := path[i-1], path[i]
		viable = viable && pf.distance(prev, next) == 1
	}
	return viable
}

// FindPath searches for a path from src to dst.
//
// When a path cannot be found it reports so and returns a path containing exactly src.
// Otherwise it returns a viable path.
func (pf PathFinder) FindPath(src, dst Position) (path Path, exists bool) {
	foundPath := astar.FindPath(pf.graph(), src, dst, pf.distance, pf.heuristic)
	if len(foundPath) == 0 {
		return Path{src}, false
	}
	path = make(Path, len(foundPath))
	for i := range path {
		path[i] = foundPath[i].(Position)
	}
	return path, true
}

func (pf PathFinder) graph() astar.Graph {
	return &_PathFinderGraph{space: pf.space}
}

func (pf PathFinder) distance(a, b astar.Node) float64 {
	h := pf.heuristic(a, b)
	if h != 1 {
		return math.Inf(0)
	}
	return 1
}

func (PathFinder) heuristic(a, b astar.Node) float64 {
	first, second := a.(Position).AtPoint(), b.(Position).AtPoint()
	return math.Abs(float64(second.Row-first.Row)) + math.Abs(float64(second.Column-first.Column))
}

type _PathFinderGraph struct {
	neighbourBuf [4]astar.Node
	space        *Space
}

var _ astar.Graph = &_PathFinderGraph{}

func (g *_PathFinderGraph) Neighbours(node astar.Node) []astar.Node {
	pt := node.(Position).AtPoint()

	ns := g.neighbourBuf[:0]
	ns = g.appendViable(ns, P(pt.Row, pt.Column+1))
	ns = g.appendViable(ns, P(pt.Row+1, pt.Column))
	ns = g.appendViable(ns, P(pt.Row, pt.Column-1))
	ns = g.appendViable(ns, P(pt.Row-1, pt.Column))
	return ns
}

func (g *_PathFinderGraph) appendViable(ns []astar.Node, pt Point) []astar.Node {
	pos := g.space.At(pt)
	if !pos.Exists() || pos.Taken() {
		return ns
	}
	return append(ns, pos)
}
