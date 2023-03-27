// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"image"
	"os"
	"sort"

	"github.com/szabba/tob-cob/ui/draw"
	"github.com/szabba/tob-cob/ui/geometry"
	"golang.org/x/exp/slog"
)

// An Anchor computes a point to treat as the position of an image, relative to image bounds.
// The computed point does not have to live within the image bounds.
type Anchor func(bounds geometry.Rect) (offset geometry.Vec)

func (anc Anchor) For(bounds geometry.Rect) (offset geometry.Vec) {
	return anc(bounds)
}

// AnchorNorthWest sets the anchor point at the upper-leftmost corner of the image.
func AnchorNorthWest() Anchor { return anchorNorthWest }

func anchorNorthWest(bounds geometry.Rect) geometry.Vec {
	// return geometry.V(0, -bounds.H())
	return geometry.V(bounds.Min.X, bounds.Max.Y)
}

// AnchorSouth sets the anchor point at the middle of the lower edge of the image.
func AnchorSouth() Anchor { return anchorSouth }

func anchorSouth(bounds geometry.Rect) geometry.Vec {
	// return geometry.V(-bounds.W()/2, 0)
	return geometry.V(
		bounds.Min.X+bounds.W()/2,
		bounds.Min.Y,
	)
}

// AnchorCenter sets the anchor point in the middle of the image.
func AnchorCenter() Anchor { return anchorCenter }

func anchorCenter(bounds geometry.Rect) geometry.Vec {
	// return geometry.V(-bounds.W()/2, -bounds.H()/2)
	return bounds.Center()
}

type Sprite struct {
	img       draw.Image
	anchor    geometry.Vec // in the image bounds coordinate system
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
	anchorPoint := anchor.For(dstImg.Bounds())
	slog.Info("sprite loaded",
		slog.String("filename", fname),
		slog.Any("bounds", dstImg.Bounds()),
		slog.Any("anchor-point", anchorPoint),
	)
	return &Sprite{dstImg, anchorPoint, geometry.Identity()}, nil
}

func (s Sprite) Draw() {
	s.img.Draw(s.matrix())
}

func (s Sprite) matrix() geometry.Mat {
	offset := s.anchor.Sub(s.img.Bounds().Min)
	t := geometry.Translation(offset)
	out := s.transform.Compose(t)

	if slog.Default().Enabled(nil, slog.LevelDebug) {

		slog.Info(
			"calculated sprite matrix",
			slog.Any("composed", out),
			slog.Any("offset-t", t),
			slog.Any("transform", s.transform),
		)
	}

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
