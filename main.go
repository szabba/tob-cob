// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/rs/zerolog/log"
)

func main() {
	pixelgl.Run(run)
}

var (
	Black = pixel.RGB(0, 0, 0)
	White = pixel.RGB(1, 1, 1)
	Red   = pixel.RGB(1, 0, 0)
)

func run() {
	wcfg := pixelgl.WindowConfig{
		Title:  "Tears of Butterflies: Colors of Blood",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}

	grid := Grid{
		CellWidth:  50,
		CellHeight: 30,
		Dx:         5,
		Dy:         5,
	}

	humanoidSprite, err := loadSprite("assets/humanoid.png")
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	w, err := pixelgl.NewWindow(wcfg)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	for !w.Closed() {
		w.Update()
		w.Clear(Black)
		center := w.Bounds().Center()
		w.Canvas().SetMatrix(pixel.IM.Scaled(pixel.ZV, 2).Moved(center))
		grid.DrawCell(w, 0, 0, White)
		grid.DrawCell(w, 0, 1, White)
		grid.DrawCell(w, 1, 0, White)
		grid.DrawCell(w, 1, 1, White)
		humanoidSprite.Move(grid, 1, 1).Draw(w)
		humanoidSprite.Move(grid, 0, 1).Draw(w)
		humanoidSprite.Move(grid, 1, 0).Draw(w)
		humanoidSprite.Move(grid, 0, 0).Draw(w)
	}
}

func loadSprite(fname string) (*Sprite, error) {
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

type Sprite struct {
	*pixel.Sprite
	offset, matrix pixel.Matrix
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

type Grid struct {
	CellWidth  float64
	CellHeight float64
	Dx, Dy     float64
}

func (grid Grid) Matrix(col, row int) pixel.Matrix {
	x, y := float64(col), float64(row)
	dx := grid.CellWidth*x + grid.Dx*x
	dy := grid.CellHeight*y + grid.Dy*y
	dr := pixel.V(dx, dy)
	return pixel.IM.Moved(dr)
}

func (grid Grid) DrawCell(dst pixel.Target, col, row int, color pixel.RGBA) {
	imd := imdraw.New(nil)
	imd.Color = color
	imd.SetMatrix(grid.Matrix(col, row))
	wHalf := float64(grid.CellWidth) / 2
	imd.Push(pixel.V(-wHalf, 0), pixel.V(wHalf, grid.CellHeight))
	imd.Rectangle(1)
	imd.Draw(dst)
}
