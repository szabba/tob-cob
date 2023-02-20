// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"image"
	"os"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/geometry"
)

type Anchor func(bounds geometry.Rect) (offset geometry.Vec)

func (anc Anchor) For(bounds geometry.Rect) (offset geometry.Vec) {
	return anc(bounds)
}

func AnchorNorthWest() Anchor { return anchorNorthWest }

func anchorNorthWest(bounds geometry.Rect) geometry.Vec {
	return geometry.V(bounds.W()/2, -bounds.H()/2)
}

func AnchorSouth() Anchor { return anchorSouth }

func anchorSouth(bounds geometry.Rect) geometry.Vec {
	return geometry.V(0, bounds.H()/2)
}

type Sprite struct {
	img draw.Image
	// offset is the vector by which the sprite has to be moved to ensure the correct anchor point.
	offset    geometry.Vec
	transform geometry.Mat
}

func LoadSprite(fname string, dst draw.Target, anchor Anchor) (*Sprite, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	dstImg := dst.Import(img)
	offset := anchor.For(dstImg.Bounds())
	return &Sprite{dstImg, offset, geometry.Identity()}, nil
}

func (s Sprite) Draw() {
	s.img.Draw(s.matrix())
}

func (s Sprite) matrix() geometry.Mat {
	t := geometry.Translation(s.offset)
	out := s.transform.Compose(t)
	log.Info().
		Str("matrix.out", out.String()).
		Str("matrix.offsetT", t.String()).
		Str("matrix.transform", s.transform.String()).
		Msg("calculated sprite matrix")
	return out
}

func (s Sprite) Transform(m geometry.Mat) Sprite {
	s.transform = m.Compose(s.transform)
	return s
}

// An OrderedSpriteGroup keeps track of a bunch of sprites and knows how to draw in the correct order.
// This assumes all the sprites are anchored at their bottom.
type OrderedSpriteGroup struct {
	order []Sprite
}

// Add sprites to draw.
// You need to add the sprites before each Draw call.
func (group *OrderedSpriteGroup) Add(sprites ...Sprite) {
	group.order = append(group.order, sprites...)
}

// Draw the added sprites.
// The set of sprites to draw and their order is forgotten afterwards.
func (group *OrderedSpriteGroup) Draw() {
	defer group.empty()
	group.sort()
	for _, sprite := range group.order {
		sprite.Draw()
	}
}

func (group *OrderedSpriteGroup) empty() {
	group.order = group.order[:0]
}

func (group *OrderedSpriteGroup) sort() {
	sort.SliceStable(group.order, func(i, j int) bool {
		first, second := group.order[i], group.order[j]
		return group.yOf(first) >= group.yOf(second)
	})
}

func (OrderedSpriteGroup) yOf(s Sprite) float64 {
	origin := geometry.V(0, 0)
	return s.transform.Apply(origin).Y
}
