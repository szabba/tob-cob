// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ui

import (
	"image"
	"image/color"
	"os"

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
	sprite.DrawOutline(dst, pixel.RGB(1, 0, 0))
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

func (sprite Sprite) DrawOutline(dst pixel.Target, color color.Color) {
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

	outline.Draw(dst)
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
