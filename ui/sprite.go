// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"image"
	"image/color"
	"os"
	"sort"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Anchor func(bounds pixel.Rect) (offset pixel.Vec)

func (anc Anchor) For(bounds pixel.Rect) (offset pixel.Vec) {
	return anc(bounds)
}

func AnchorNorthWest() Anchor { return anchorNorthWest }

func anchorNorthWest(bounds pixel.Rect) pixel.Vec {
	return pixel.V(bounds.W()/2, -bounds.H()/2)
}

func AnchorSouth() Anchor { return anchorSouth }

func anchorSouth(bounds pixel.Rect) pixel.Vec {
	return pixel.V(0, bounds.H()/2)
}

type Sprite struct {
	*pixel.Sprite
	// offset is the vector by which the sprite has to be moved to ensure the correct anchor point.
	offset    pixel.Vec
	transform pixel.Matrix
}

func LoadSprite(fname string, anchor Anchor) (*Sprite, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	offset := anchor.For(pic.Bounds())
	return &Sprite{sprite, offset, pixel.IM}, nil
}

func (sprite Sprite) Draw(dst pixel.Target) {
	sprite.Sprite.Draw(dst, sprite.matrix())
}

func (sprite Sprite) matrix() pixel.Matrix {
	m := pixel.IM
	m = m.Moved(sprite.offset)
	m = m.Chained(sprite.transform)
	return m
}

func (sprite Sprite) Transform(m pixel.Matrix) Sprite {
	sprite.transform = sprite.transform.Chained(m)
	return sprite
}

func (sprite Sprite) Outline(color color.Color, width float64) Outline {
	frame := sprite.Sprite.Frame()
	frame = frame.Moved(frame.Center().Scaled(-1))
	frame = frame.Moved(sprite.offset)
	frame.Min = sprite.transform.Project(frame.Min)
	frame.Max = sprite.transform.Project(frame.Max)
	outline := Outline{
		Rect:  frame,
		Color: color,
		Width: 1,
	}
	return outline
}

type Outline struct {
	Color color.Color
	Width float64
	pixel.Rect
}

func (out Outline) Draw(dst pixel.Target) {
	imd := imdraw.New(nil)
	imd.Color = out.Color

	imd.Push(out.Min, out.Max)
	imd.Rectangle(1)
	imd.Draw(dst)

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
func (group *OrderedSpriteGroup) Draw(dst pixel.Target) {
	defer group.empty()
	group.sort()
	for _, sprite := range group.order {
		sprite.Draw(dst)
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

func (OrderedSpriteGroup) yOf(sprite Sprite) float64 {
	origin := pixel.V(0, 0)
	return sprite.transform.Project(origin).Y
}
