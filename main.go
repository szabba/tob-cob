// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/rs/zerolog/log"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	wcfg := pixelgl.WindowConfig{
		Title:  "Tears of Butterflies: Colors of Blood",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
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
	grid := Grid{}
	grid.Cell.Width = 50
	grid.Cell.Heigth = 25
	grid.Spacing.Dx = 5
	grid.Spacing.Dy = 5

	for !w.Closed() {
		w.Update()
		center := w.Bounds().Center()
		w.Canvas().SetMatrix(pixel.IM.Scaled(center.Scaled(-1), 2))
		humanoidSprite.Draw(w, grid.Matrix(1, 1))
		humanoidSprite.Draw(w, grid.Matrix(0, 1))
		humanoidSprite.Draw(w, grid.Matrix(1, 0))
		humanoidSprite.Draw(w, grid.Matrix(0, 0))
	}
}

func loadSprite(fname string) (*pixel.Sprite, error) {
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
	return sprite, nil
}

type Grid struct {
	Cell struct {
		Width, Heigth float64
	}
	Spacing struct {
		Dx, Dy float64
	}
}

func (grid Grid) Matrix(col, row int) pixel.Matrix {
	x, y := float64(col), float64(row)
	dx := grid.Cell.Width*x + grid.Spacing.Dx*x
	dy := grid.Cell.Heigth*y + grid.Spacing.Dy*y
	dr := pixel.V(dx, dy)
	return pixel.IM.Moved(dr)
}
