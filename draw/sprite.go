// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package draw

import (
	"image"
	"image/color"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Sprite struct {
	*pixel.Sprite
	offset, matrix pixel.Matrix
}

func LoadSprite(fname string) (*Sprite, error) {
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
	offset := pixel.V(0, pic.Bounds().H()/2)
	offsetMatrix := pixel.IM.Moved(offset)
	return &Sprite{sprite, offsetMatrix, pixel.IM}, nil
}

func (sprite Sprite) Draw(dst pixel.Target) {
	sprite.DrawOutline(dst, Red)
	sprite.Sprite.Draw(dst, sprite.offset.Chained(sprite.matrix))
}

func (sprite Sprite) Move(grid Grid, col, row int) Sprite {
	sprite.matrix = sprite.matrix.Chained(grid.Matrix(col, row))
	return sprite
}

func (sprite Sprite) DrawOutline(dst pixel.Target, color color.Color) {
	frame := sprite.Sprite.Frame()
	wHalf := frame.W() / 2
	outlineOffset := pixel.V(-wHalf, 0)
	outline := frame.Moved(outlineOffset)

	imd := imdraw.New(nil)
	imd.Color = color
	imd.SetMatrix(sprite.matrix)

	imd.Push(outline.Min, outline.Max)
	imd.Rectangle(1)
	imd.Draw(dst)
}
